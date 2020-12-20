package api

import (
	"ditto/booking/config"
	"ditto/booking/mail"
	"ditto/booking/utils"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreateClass - クラスの仮作成
// @Summary クラスの仮作成を行います
// @Tags Class
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /class [post]
func (s *Service) CreateClass(c echo.Context) error {
	logon := s.logonFromToken(c)

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//input check
	//name
	if _, ok := input["name"]; !ok {
		return BadRequest(errors.New("Class name is required"))
	}

	tx := s.DB().Begin()
	//create class
	class, err := s.DB().CreateClass(tx, logon, input)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}
	tx.Commit()

	resp := Response{
		Code: http.StatusOK,
		Data: class,
	}

	return c.JSON(http.StatusOK, resp)
}

// ChangeUserClass - ユーザークラスを切り替え
// @Summary ユーザークラスを切り替えします
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "class id"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /user/class/{id} [put]
func (s *Service) ChangeUserClass(c echo.Context) error {
	logon := s.logonFromToken(c)

	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Id is required"))
	}
	tid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}
	tx := s.DB().Begin()

	err = s.DB().ChangeUserClass(tx, logon, tid)
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

// GetClasses - ユーザーのクラス一覧取得
// @Summary ユーザーのクラス一覧取得を取得します
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /user/classes [get]
func (s *Service) GetClasses(c echo.Context) error {
	logon := s.logonFromToken(c)

	classes, err := s.DB().GetUserClasses(nil, logon, logon.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return NotFound(err)
		}
		return InternalServerError(err)
	}
	if len(classes) == 0 {
		return NotFound(errors.New("No content"))
	}

	resp := Response{
		Code: http.StatusOK,
		Data: classes,
	}

	return c.JSON(http.StatusOK, resp)
}

// ClassListUserWithDetail - ユーザー一覧取得(詳細)
// @Summary ユーザー一覧（詳細）を取得します
// @Tags Class
// @Accept json
// @Produce json
// @Param id path int true "class id"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /class/users/{id}/detail [get]
func (s *Service) ClassListUserWithDetail(c echo.Context) error {
	logon := s.logonFromToken(c)

	if logon.Tenant == 0 {
		return NoContent(errors.New("No tenant"))
	}

	//class id
	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Class id is required"))
	}
	cid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	users, err := s.DB().GetClassUsersWithDetail(nil, logon, cid)
	if err != nil {
		return InternalServerError(err)
	}
	if len(users) == 0 {
		return NotFound(errors.New("No content"))
	}
	resp := Response{
		Code: http.StatusOK,
		Data: users,
	}

	return c.JSON(http.StatusOK, resp)
}

// ClassListUser - ユーザー一覧取得
// @Summary クラスのユーザー一覧を取得します
// @Tags Class
// @Accept json
// @Produce json
// @Param id path int true "class id"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /class/users/{id} [get]
func (s *Service) ClassListUser(c echo.Context) error {
	logon := s.logonFromToken(c)

	if logon.Tenant == 0 {
		return NoContent(errors.New("No tenant"))
	}

	//class id
	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Class id is required"))
	}
	cid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	users, err := s.DB().GetClassUsers(nil, logon, cid)
	if err != nil {
		return InternalServerError(err)
	}
	if len(users) == 0 {
		return NotFound(errors.New("No content"))
	}
	resp := Response{
		Code: http.StatusOK,
		Data: users,
	}

	return c.JSON(http.StatusOK, resp)
}

// ClassCreateUser - ユーザー作成
// @Summary ユーザーを作成します（クラス）
// @Tags Class
// @Accept json
// @Produce json
// @Param id path int true "class id"
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /class/user/{id} [post]
func (s *Service) ClassCreateUser(c echo.Context) error {
	logon := s.logonFromToken(c)

	if logon.Tenant == 0 {
		return NoContent(errors.New("No tenant"))
	}

	//class id
	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Class id is required"))
	}
	cid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	conf := config.Load()

	//input check
	if _, ok := input["email"]; !ok {
		return BadRequest(errors.New("Email is required"))
	}
	//password
	if _, ok := input["password"]; !ok {
		pass := utils.GeneratePassowrd(8, false)
		input["password"] = pass
	}
	password := input["password"].(string)

	tx := s.DB().Begin()
	account, err := s.DB().ClassCreateUser(tx, logon, cid, input)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	//user details
	err = s.DB().UpdateUser(tx, logon, account.ID, input)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	role := conf.DefaultRole()
	err = s.DB().AddUserRole(tx, logon, account.ID, []int64{role.ID})
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

	//ge class
	class, err := s.DB().GetClass(tx, logon, cid, "")
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	email := account.Email
	confirmurl := fmt.Sprintf(conf.Confirm.URL, email, confirm.ConfirmCode)
	confirmurl = url.QueryEscape(confirmurl)
	val := map[string]interface{}{
		"LessonName": class.Name,
		"Email":      email,
		"Expire":     conf.Confirm.Expires,
		"ConfirmURL": confirmurl,
		"Password":   password,
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

	//clear password
	account.PasswordHash = ""
	account.TemporaryPassword = password

	resp := Response{
		Code: http.StatusOK,
		Data: account,
	}

	return c.JSON(http.StatusOK, resp)
}

// ClassDividedUser - ユーザーの所属
// @Summary ユーザーとテナントの所属を変更します（クラス）
// @Tags Class
// @Accept json
// @Produce json
// @Param id path int true "class id"
// @Param users body DivideUsers true "user list"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /class/user/{id} [put]
func (s *Service) ClassDividedUser(c echo.Context) error {
	logon := s.logonFromToken(c)

	//テナント
	if logon.Tenant == 0 {
		return BadRequest(errors.New("Logon user hasn't a valid tenant"))
	}

	//class id
	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Class id is required"))
	}
	cid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	data := DivideUsers{}
	//decode
	if err := c.Bind(&data); err != nil {
		return err
	}

	au := make([]int64, 0) //add
	ru := make([]int64, 0) //remove
	//add
	for _, u := range data.List {
		if u.Divides == 0 {
			ru = append(ru, u.UserID)
		} else if u.Divides == 1 {
			au = append(au, u.UserID)
		}
	}
	tx := s.DB().Begin()
	//remove
	if len(ru) > 0 {
		err := s.DB().RemoveUserFromClass(tx, logon, cid, ru)
		if err != nil {
			tx.Rollback()
			return InternalServerError(err)
		}
	}
	//add
	if len(au) > 0 {
		err := s.DB().DivideUserToClass(tx, logon, cid, au)
		if err != nil {
			tx.Rollback()
			return InternalServerError(err)
		}
	}
	tx.Commit()

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}
