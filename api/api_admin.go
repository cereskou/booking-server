package api

import (
	"ditto/booking/config"
	"ditto/booking/logger"
	"ditto/booking/mail"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// AdminGetUser - ユーザー情報
// @Summary ユーザー情報を取得します(admin)
// @Tags Admin/user
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/user/{id} [get]
func (s *Service) AdminGetUser(c echo.Context) error {
	//
	uid, err := paramInt(c, "id", "User id is required")
	if err != nil {
		return err
	}

	user, err := s.DB().GetUser(nil, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return NotFound(err)
		}
		return InternalServerError(err)
	}

	resp := Response{
		Code: http.StatusOK,
		Data: user,
	}

	return c.JSON(http.StatusOK, resp)
}

// AdminUpdateUser - ユーザー情報を更新します
// @Summary ユーザー情報を更新します(admin)
// @Tags Admin/user
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/user/{id} [put]
func (s *Service) AdminUpdateUser(c echo.Context) error {
	logon := s.logonFromToken(c)

	uid, err := paramInt(c, "id", "User id is required")
	if err != nil {
		return err
	}

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	tx := s.DB().Begin()
	//update
	err = s.DB().UpdateUser(tx, logon, uid, input)
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

// AdminCreateAccount - アカウント情報作成
// @Summary アカウント情報を新規作成します(admin)
// @Tags Admin/user
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/user [post]
func (s *Service) AdminCreateAccount(c echo.Context) error {
	logon := s.logonFromToken(c)

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
	account, err := s.DB().CreateAccount(tx, logon, input)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	//Create a confirm code
	confirm, err := s.DB().CreateConfirmCode(tx, account)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	//send mail
	logger.Trace(confirm.ConfirmCode)

	//update
	err = s.DB().UpdateUser(tx, logon, account.ID, input)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	conf := config.Load()

	email := account.Email
	confirmurl := fmt.Sprintf(conf.Confirm.URL, email, confirm.ConfirmCode)
	confirmurl = url.QueryEscape(confirmurl)

	val := map[string]interface{}{
		"LessonName": "Lesson",
		"Email":      email,
		"Expire":     conf.Confirm.Expires,
		"ConfirmURL": confirmurl,
	}
	mt, err := s.DB().GetMailTemplate(tx, 0, "mailconfirm")
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}
	msend := mail.New()
	body, err := msend.Render(mt.MailID, mt.Body, val)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}
	err = msend.Send(email, email, mt.Subject, body)
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

// AdminDeleteAcount - アカウント削除
// @Summary アカウントを削除します(admin)
// @Tags Admin/user
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/user/{id} [delete]]
func (s *Service) AdminDeleteAcount(c echo.Context) error {
	logon := s.logonFromToken(c)

	uid, err := paramInt(c, "id", "User id is required")
	if err != nil {
		return err
	}

	tx := s.DB().Begin()

	//1. user delete
	err = s.DB().DeleteUser(tx, uid)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}
	//2. delete account
	err = s.DB().DeleteAccount(tx, logon, uid)
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

// AdminGetAccount - アカウント情報取得
// @Summary アカウント情報取得します(admin)
// @Tags Admin/user
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/user/{id}/account [get]
func (s *Service) AdminGetAccount(c echo.Context) error {
	uid, err := paramInt(c, "id", "User id is required")
	if err != nil {
		return err
	}

	user, err := s.DB().GetAccountByID(nil, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return NotFound(err)
		}
		return InternalServerError(err)
	}
	user.PasswordHash = ""

	resp := Response{
		Code: http.StatusOK,
		Data: user,
	}

	return c.JSON(http.StatusOK, resp)
}

// AdminCreateTenant - テナント作成
// @Summary テナントを新規作成します(admin)
// @Tags Admin/tenant
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/tenant [post]
func (s *Service) AdminCreateTenant(c echo.Context) error {
	logon := s.logonFromToken(c)

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//input check
	if _, ok := input["name"]; !ok {
		return BadRequest(errors.New("Name is required"))
	}

	tx := s.DB().Begin()

	//create tenant and return
	tenant, err := s.DB().CreateTenant(tx, logon, input)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}
	tx.Commit()

	resp := Response{
		Code: http.StatusOK,
		Data: tenant,
	}

	return c.JSON(http.StatusOK, resp)
}

// AdminGetTenant - テナント情報取得
// @Summary テナント情報を取得します(admin)
// @Tags Admin/tenant
// @Accept json
// @Produce json
// @Param id path int true "tenant id"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/tenant/{id} [get]
func (s *Service) AdminGetTenant(c echo.Context) error {
	logon := s.logonFromToken(c)

	id, err := paramInt(c, "id", "Tenant id is required")
	if err != nil {
		return err
	}

	//update
	tenant, err := s.DB().GetTenant(nil, logon, id, "")
	if err != nil {
		return InternalServerError(err)
	}

	resp := Response{
		Code: http.StatusOK,
		Data: tenant,
	}

	return c.JSON(http.StatusOK, resp)
}

// AdminUpdateTenant - テナント情報を更新します
// @Summary テナント情報を更新します(admin)
// @Tags Admin/tenant
// @Accept json
// @Produce json
// @Param id path string true "tenant id"
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/tenant/{id} [put]
func (s *Service) AdminUpdateTenant(c echo.Context) error {
	logon := s.logonFromToken(c)

	id, err := paramInt(c, "id", "Tenant id is required")
	if err != nil {
		return err
	}

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//input check
	if _, ok := input["name"]; !ok {
		return BadRequest(errors.New("Name is required"))
	}

	tx := s.DB().Begin()
	//update
	err = s.DB().UpdateTenant(tx, logon, id, input)
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

// AdminDeleteTenant - テナント情報削除
// @Summary テナント情報を削除します(admin)
// @Tags Admin/tenant
// @Accept json
// @Produce json
// @Param id path string true "tenant id"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/tenant/{id} [delete]]
func (s *Service) AdminDeleteTenant(c echo.Context) error {
	logon := s.logonFromToken(c)

	id, err := paramInt(c, "id", "Tenant id is required")
	if err != nil {
		return err
	}

	tx := s.DB().Begin()

	//1. tentans delete
	err = s.DB().DeleteTenant(tx, logon, id)
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
