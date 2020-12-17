package api

import (
	"ditto/booking/config"
	"ditto/booking/logger"
	"ditto/booking/mail"
	"fmt"
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

	user, err := s.DB().GetUser(nil, email)
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

	user, err := s.DB().GetAccount(nil, email)
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

	tx := s.DB().Begin()
	//update
	err := s.DB().UpdateUser(tx, logon.ID, email, input)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

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
	tx := s.DB().Begin()

	//update
	account, err := s.DB().CreateAccount(tx, logon.ID, input)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Create a confirm code
	confirm, err := s.DB().CreateConfirmCode(tx, account)
	if err != nil {
		tx.Rollback()
		return err
	}

	//send mail
	logger.Trace(confirm.ConfirmCode)

	email := account.Email
	//update
	err = s.DB().UpdateUser(tx, logon.ID, email, input)
	if err != nil {
		tx.Rollback()
		return err
	}

	conf := config.Load()

	url := fmt.Sprintf(conf.Confirm.URL, email, confirm.ConfirmCode)
	val := map[string]interface{}{
		"LessonName": "Lesson",
		"Email":      email,
		"Expire":     conf.Confirm.Expires,
		"ConfirmURL": url,
	}
	mt, err := s.DB().GetMailTemplate(tx, 0, "mailconfirm")
	if err != nil {
		tx.Rollback()
		return err
	}
	msend := mail.New()
	body, err := msend.Render(mt.MailID, mt.Body, val)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = msend.Send(email, email, mt.Subject, body)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}

// AdminDeleteAcount - アカウント削除
// @Summary アカウントを削除します(admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param email path string true "email"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/user/{email} [delete]]
func (s *Service) AdminDeleteAcount(c echo.Context) error {
	logon := logonFromToken(c)
	email := c.Param("email")
	email, _ = url.QueryUnescape(email)

	tx := s.DB().Begin()

	//1. user delete
	err := s.DB().DeleteUser(tx, email)
	if err != nil {
		tx.Rollback()
		return err
	}
	//2. delete account
	err = s.DB().DeleteAccount(tx, logon.ID, email)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}
