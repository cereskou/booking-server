package api

import (
	"ditto/booking/config"
	"ditto/booking/cx"
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

// Logout - ログアウト
// @Summary ログアウト
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /user/logout [get]
func (s *Service) Logout(c echo.Context) error {
	logon := s.logonFromToken(c)

	key := "ACCN_" + logon.Email
	err := s.CacheDel(key)
	if err != nil {
		return InternalServerError(err)
	}
	key = "ACCN_TENANT_" + logon.Email
	err = s.CacheDel(key)
	if err != nil {
		return InternalServerError(err)
	}

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}

// Login - ログイン
// @Summary ログイン
// @Tags Login
// @Accept json
// @Produce json
// @Param data body Login false "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Router /login [post]
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
			return BadRequest(errors.New("account not confirmed"))
		}
		//compare password
		if !utils.CompareHashedPassword(u.PasswordHash, login.Password) {
			return c.JSON(http.StatusUnauthorized, Unauthorized())
		}

		err = copier.Copy(&user, u)
		if err != nil {
			return InternalServerError(err)
		}

		key := "ACCN_" + login.Email
		err = s.CacheSet(key, user)
		if err != nil {
		}
		//set tenat
		tenant, err := s.DB().GetUserTenant(nil, u.ID)
		if err != nil {
			//
		} else {
			key = "ACCN_TENANT_" + login.Email
			err = s.CacheSet(key, tenant)
			if err != nil {
			}
		}
	} else {
		//compare password
		if !utils.CompareHashedPassword(user.PasswordHash, login.Password) {
			return c.JSON(http.StatusUnauthorized, Unauthorized())
		}

	}

	//payload
	d := cx.Payload{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	logger.Trace(d)

	//create a token
	token, err := s.generateToken(&d)
	if err != nil {
		return InternalServerError(err)
	}
	resp := Response{
		Code: http.StatusOK,
		Data: token,
	}

	return c.JSON(http.StatusOK, resp)
}

// RefreshToken - api to refresh tokens
// @Summary リフレッシュトークン
// @Tags Login
// @Accept json
// @Produce json
// @Param data body RefreshToken false "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /refresh [post]
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

		var d cx.Payload
		err := utils.JSON.NewDecoder(strings.NewReader(payload)).Decode(&d)
		if err != nil {
			return InternalServerError(err)
		}

		//create a token
		token, err := s.generateToken(&d)
		if err != nil {
			return InternalServerError(err)
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
// @Tags Login
// @Accept json
// @Produce json
// @Param code query string true "are確認コード"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Router /confirm [get]
func (s *Service) ConfirmEmail(c echo.Context) error {
	code, err := queryParam(c, "code", "Code is required")
	if err != nil {
		return err
	}

	conf := config.Load()

	//有効期間
	expires := utils.HourToSecond(conf.Confirm.Expires)

	//begin transaction
	tx := s.DB().Begin()

	//get confirm record
	err = s.DB().ConfirmAccountWithCode(tx, code, expires)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}
	tx.Commit()

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}
