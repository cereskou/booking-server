package api

import (
	"ditto/booking/logger"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

// AdminGetUser - ユーザー情報
// @Summary ユーザー情報を取得します(admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param email path string true "email"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/user/{email} [get]
func (s *Service) AdminGetUser(c echo.Context) error {
	email := c.Param("email")
	email, _ = url.QueryUnescape(email)

	user, err := s.DB().GetUser(email)
	if err != nil {
		resp := Response{
			Code:  http.StatusNotFound,
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotFound, resp)
	}

	resp := Response{
		Code: http.StatusOK,
		Data: user,
	}

	return c.JSON(http.StatusOK, resp)
}

// AdminGetAccount - アカウント情報取得
// @Summary アカウント情報取得します(admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param email path string true "email"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/account/{email} [get]
func (s *Service) AdminGetAccount(c echo.Context) error {
	email := c.Param("email")
	email, _ = url.QueryUnescape(email)

	user, err := s.DB().GetAccount(email)
	if err != nil {
		resp := Response{
			Code:  http.StatusNotFound,
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotFound, resp)
	}
	user.PasswordHash = ""

	resp := Response{
		Code: http.StatusOK,
		Data: user,
	}

	return c.JSON(http.StatusOK, resp)
}

// AdminUpdateUser - ユーザー情報を更新します
// @Summary ユーザー情報を更新します(admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/user/{email} [put]
func (s *Service) AdminUpdateUser(c echo.Context) error {
	logon := logonFromToken(c)
	email := c.Param("email")
	email, _ = url.QueryUnescape(email)

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//update
	err := s.DB().UpdateUser(logon.ID, email, input)
	if err != nil {
		return err
	}

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}

// AdminCreateAccount - アカウント情報作成
// @Summary アカウント情報を新規作成します(admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/user [post]
func (s *Service) AdminCreateAccount(c echo.Context) error {
	logon := logonFromToken(c)

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}
	// input check
	//* email
	//* password
	// role
	//* name
	// age
	// phone
	// contact
	// gender
	// occupation

	//update
	account, err := s.DB().CreateAccount(logon.ID, input)
	if err != nil {
		return err
	}

	//Create a confirm code
	confirm, err := s.DB().CreateConfirmCode(account)
	if err != nil {
		return err
	}
	//send mail
	logger.Trace(confirm.ConfirmCode)

	email := account.Email
	//update
	err = s.DB().UpdateUser(logon.ID, email, input)
	if err != nil {
		return err
	}

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}
