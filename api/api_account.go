package api

import (
	"ditto/booking/logger"
	"ditto/booking/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Login - ログイン
// @Summary ログイン
// @Tags Account
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

	var role string
	var email string
	var name string
	//check cache
	info, err := s.HGetAll(login.Email)
	if len(info) == 0 || err != nil {
		logger.Trace("Find user in db")
		//
		//get user from db
		user, err := s.DB().GetUser(login.Email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, Unauthorized())
		}

		//compare password
		if !utils.CompareHashedPassword(user.PasswordHash, login.Password) {
			return c.JSON(http.StatusUnauthorized, Unauthorized())
		}

		//set user role
		err = s.HSet(login.Email, "role", user.Role)
		if err != nil {
		}
		//set user email
		err = s.HSet(login.Email, "email", user.Email)
		if err != nil {
		}
		err = s.HSet(login.Email, "name", user.Name)
		if err != nil {
		}

		role = user.Role
		email = user.Email
		name = user.Name
	} else {
		logger.Trace("Found in cache")

		role = info["role"]
		email = info["email"]
		name = info["name"]
	}

	//create a token
	token, err := s.generateToken(name, email, role)
	if err != nil {
		return err
	}
	resp := Response{
		Code: http.StatusOK,
		Data: token,
	}

	return c.JSON(http.StatusOK, resp)
}
