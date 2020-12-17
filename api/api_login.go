package api

import (
	"ditto/booking/config"
	"ditto/booking/logger"
	"ditto/booking/models"
	"ditto/booking/security"
	"ditto/booking/utils"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

// Login - ログイン
// @Summary ログイン
// @Tags User
// @Accept json
// @Produce json
// @Param data body Login false "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Router /user/login [post]
func (s *Service) Login(c echo.Context) error {
	login := Login{}
	//decode
	if err := c.Bind(&login); err != nil {
		return err
	}

	var user models.AccountWithRole
	//check cache
	err := s.CacheGet(login.Email, &user)
	if err != nil {
		logger.Trace("Find user in db")
		//
		//get user from db
		u, err := s.DB().GetAccount(nil, login.Email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, Unauthorized())
		}
		//confirmed
		if u.EmailConfirmed == 0 {
			return c.JSON(http.StatusBadRequest, BadRequest(errors.New("account not confirmed")))
		}
		//compare password
		if !utils.CompareHashedPassword(u.PasswordHash, login.Password) {
			return c.JSON(http.StatusUnauthorized, Unauthorized())
		}

		err = copier.Copy(&user, u)
		if err != nil {
			return err
		}

		err = s.CacheSet(login.Email, user)
		if err != nil {
		}
	} else {
		//compare password
		if !utils.CompareHashedPassword(user.PasswordHash, login.Password) {
			return c.JSON(http.StatusUnauthorized, Unauthorized())
		}

	}

	//payload
	d := Payload{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}
	//create a token
	token, err := s.generateToken(&d)
	if err != nil {
		return err
	}
	resp := Response{
		Code: http.StatusOK,
		Data: token,
	}

	return c.JSON(http.StatusOK, resp)
}

// RefreshToken - api to refresh tokens
// @Summary リフレッシュトークン
// @Tags User
// @Accept json
// @Produce json
// @Param data body RefreshToken false "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Router /user/refresh [post]
func (s *Service) RefreshToken(c echo.Context) error {
	req := RefreshToken{}

	//decode
	if err := c.Bind(&req); err != nil {
		return err
	}

	//refresh token
	if req.GrantType != "refresh_token" {
		return echo.ErrUnauthorized
	}

	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return s._rsa.GetPublicKey(), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		secret := claims["sub"].(string)
		payload := security.DecryptString(secret)

		var d Payload
		err := utils.JSON.NewDecoder(strings.NewReader(payload)).Decode(&d)
		if err != nil {
			return err
		}

		//create a token
		token, err := s.generateToken(&d)
		if err != nil {
			return err
		}
		resp := Response{
			Code: http.StatusOK,
			Data: token,
		}

		return c.JSON(http.StatusOK, resp)
	}

	return err
}

// ConfirmEmail - Email確認
// @Summary Email確認
// @Tags User
// @Accept json
// @Produce json
// @Param e query string true "Email"
// @Param code query string true "are確認コード"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Router /user/confirm [get]
func (s *Service) ConfirmEmail(c echo.Context) error {
	email := c.QueryParam("e")
	code := c.QueryParam("code")

	conf := config.Load()

	//有効期間
	expires := utils.HourToSecond(conf.Confirm.Expires)
	tx := s.DB().Begin()

	//get confirm record
	err := s.DB().ConfirmAccountWithCode(tx, email, code, expires)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusNotFound, NewResponse(http.StatusNotFound, err.Error()))
	}
	tx.Commit()

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}
