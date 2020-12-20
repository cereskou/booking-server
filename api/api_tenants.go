package api

import (
	"ditto/booking/config"
	"ditto/booking/mail"
	"ditto/booking/models"
	"ditto/booking/utils"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreateTenant - テナントの仮作成
// @Summary テナントの仮作成を行います
// @Tags Tenant
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /tenant [post]]
func (s *Service) CreateTenant(c echo.Context) error {
	logon := s.logonFromToken(c)

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//input check
	//name
	if _, ok := input["name"]; !ok {
		return BadRequest(errors.New("Tenant name is required"))
	}
	active := 0
	//active
	if val, ok := input["active"]; !ok {
		active = 0
	} else {
		if val.(bool) {
			active = 1
		}
	}

	tx := s.DB().Begin()
	//create tenant and return
	tenant, err := s.DB().CreateTenant(tx, logon, input)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	//tenant id
	logon.Tenant = tenant.ID
	key := "ACCN_TENANT_" + logon.Email
	err = s.CacheSet(key, tenant.ID)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	//add logon user to tenant/ need re-login
	users := make([]*models.TenantUsers, 0)
	users = append(users, &models.TenantUsers{
		TenantID:   tenant.ID,
		UserID:     logon.ID,
		Right:      active,
		UpdateUser: logon.ID,
	})
	err = s.DB().AddUserToTenant(tx, logon, users)
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

// ChangeUserTenant - ユーザーの所属テナントアクティブ
// @Summary ユーザーの所属テナントをアクティブします
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "tenant id"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /user/tenants/{id} [put]
func (s *Service) ChangeUserTenant(c echo.Context) error {
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

	err = s.DB().ChangeUserTenant(tx, logon, tid)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}
	tx.Commit()

	//clear
	key := "ACCN_TENANT_" + logon.Email
	err = s.CacheSet(key, tid)
	if err != nil {
	}

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}

// GetTenants - ユーザーの所属テナント一覧取得
// @Summary ユーザーの所属テナント一覧を取得します
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /user/tenants [get]
func (s *Service) GetTenants(c echo.Context) error {
	logon := s.logonFromToken(c)

	tenants, err := s.DB().GetUserTenants(nil, logon, logon.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return NotFound(err)
		}
		return InternalServerError(err)
	}
	if len(tenants) == 0 {
		return NotFound(errors.New("No content"))
	}

	resp := Response{
		Code: http.StatusOK,
		Data: tenants,
	}

	return c.JSON(http.StatusOK, resp)
}

// TenantListUserWithDetail - ユーザー一覧取得(詳細)
// @Summary ユーザー一覧（詳細）を取得します
// @Tags Tenant
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /tenant/users/detail [get]
func (s *Service) TenantListUserWithDetail(c echo.Context) error {
	logon := s.logonFromToken(c)

	if logon.Tenant == 0 {
		return NoContent(errors.New("No tenant"))
	}

	users, err := s.DB().GetTenantUserWithDetail(nil, logon, logon.Tenant)
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

// TenantListUser - ユーザー一覧取得
// @Summary ユーザー一覧を取得します（ログイン中）
// @Tags Tenant
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /tenant/users [get]
func (s *Service) TenantListUser(c echo.Context) error {
	logon := s.logonFromToken(c)

	if logon.Tenant == 0 {
		return NoContent(errors.New("No tenant"))
	}

	users, err := s.DB().GetTenantUser(nil, logon, logon.Tenant)
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

// TenantCreateUser - ユーザー作成
// @Summary ユーザーを作成します（テナント）
// @Tags Tenant
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /tenant/user [post]
func (s *Service) TenantCreateUser(c echo.Context) error {
	logon := s.logonFromToken(c)

	if logon.Tenant == 0 {
		return NoContent(errors.New("No tenant"))
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
	account, err := s.DB().TenantCreateUser(tx, logon, input)
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

	//tentans_users
	tenant, err := s.DB().GetTenant(tx, logon, logon.Tenant, "")
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

	email := account.Email
	confirmurl := fmt.Sprintf(conf.Confirm.URL, email, confirm.ConfirmCode)
	confirmurl = url.QueryEscape(confirmurl)
	val := map[string]interface{}{
		"LessonName": tenant.Name,
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

	resp := Response{
		Code: http.StatusOK,
		Data: account,
	}

	return c.JSON(http.StatusOK, resp)
}

// TenantDeleteUser - ユーザー削除
// @Summary ユーザーを削除します（テナント）
// @Tags Tenant
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /tenant/user/{id} [delete]
func (s *Service) TenantDeleteUser(c echo.Context) error {
	logon := s.logonFromToken(c)

	//user id
	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("User id is required"))
	}
	uid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	tx := s.DB().Begin()

	//1. role
	err = s.DB().DeleteUserRole(tx, logon, logon.Tenant, []int64{uid})
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	//2. tenant_users
	err = s.DB().RemoveUserFromTenant(tx, logon, logon.Tenant, []int64{uid})
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	//3. user delete
	err = s.DB().DeleteUser(tx, uid)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}

	//4. delete account
	err = s.DB().DeleteAccount(tx, logon, uid)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}
	tx.Commit()

	//自分を削除すること
	if uid == logon.ID {
		key := "ACCN_TENANT_" + logon.Email
		err = s.CacheDel(key)
		if err != nil {
		}
	}

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}

// TenantDividedUser - ユーザーの所属
// @Summary ユーザーとテナントの所属を変更します（テナント）
// @Tags Tenant
// @Accept json
// @Produce json
// @Param users body DivideUsers true "user list"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /tenant/user [put]
func (s *Service) TenantDividedUser(c echo.Context) error {
	logon := s.logonFromToken(c)

	//テナント
	if logon.Tenant == 0 {
		return BadRequest(errors.New("Logon user hasn't a valid tenant"))
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
		err := s.DB().RemoveUserFromTenant(tx, logon, logon.Tenant, ru)
		if err != nil {
			tx.Rollback()
			return InternalServerError(err)
		}
	}
	//add
	if len(au) > 0 {
		err := s.DB().DivideUserToTenant(tx, logon, logon.Tenant, au)
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
