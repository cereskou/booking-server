package api

import (
	"ditto/booking/config"
	"ditto/booking/cx"
	"ditto/booking/logger"
	"ditto/booking/mail"
	"ditto/booking/utils"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

// Signup - アカウント作成
// @Summary アカウント情報を新規作成します
// @Tags Login
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Router /signup [post]
func (s *Service) Signup(c echo.Context) error {
	data := Signup{}
	//decode
	if err := c.Bind(&data); err != nil {
		return err
	}

	//convert input to map
	input := utils.StructToJSONTagMap(data)

	tx := s.DB().Begin()
	logon := &cx.Payload{}
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
		"LessonName": conf.Service,
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
