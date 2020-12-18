package api

import (
	"ditto/booking/config"
	"ditto/booking/logger"
	"ditto/booking/mail"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
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
	sid := c.Param("id")
	uid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}

	user, err := s.DB().GetUser(nil, uid)
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
	logon := logonFromToken(c)

	sid := c.Param("id")
	uid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
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
	account, err := s.DB().CreateAccount(tx, logon, input)
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

	//update
	err = s.DB().UpdateUser(tx, logon, account.ID, input)
	if err != nil {
		tx.Rollback()
		return err
	}

	conf := config.Load()

	email := account.Email
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
	logon := logonFromToken(c)

	sid := c.Param("id")
	uid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}

	tx := s.DB().Begin()

	//1. user delete
	err = s.DB().DeleteUser(tx, uid)
	if err != nil {
		tx.Rollback()
		return err
	}
	//2. delete account
	err = s.DB().DeleteAccount(tx, logon, uid)
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
	sid := c.Param("id")
	uid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}

	user, err := s.DB().GetAccountByID(nil, uid)
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
	logon := logonFromToken(c)

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//input check
	if _, ok := input["name"]; !ok {
		return c.JSON(http.StatusBadRequest, BadRequest(errors.New("Name is required")))
	}

	tx := s.DB().Begin()

	//create tenant and return
	tenant, err := s.DB().CreateTenant(tx, logon, input)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, InternalServerError(err))
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
// @Param id query int true "tenant id"
// @Param name query string false "名前"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/tenant [get]
func (s *Service) AdminGetTenant(c echo.Context) error {
	logon := logonFromToken(c)

	sid := c.QueryParam("id")
	name := c.QueryParam("name")
	name, _ = url.QueryUnescape(name)

	tid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}

	//update
	tenant, err := s.DB().GetTenant(nil, logon, tid, name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerError(err))
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
	logon := logonFromToken(c)
	sid := c.Param("id")
	tid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//input check
	if _, ok := input["name"]; !ok {
		return c.JSON(http.StatusBadRequest, BadRequest(errors.New("Name is required")))
	}

	tx := s.DB().Begin()
	//update
	err = s.DB().UpdateTenant(tx, logon, tid, input)
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
	logon := logonFromToken(c)
	sid := c.Param("id")
	tid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}

	tx := s.DB().Begin()

	//1. tentans delete
	err = s.DB().DeleteTenant(tx, logon, tid)
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
