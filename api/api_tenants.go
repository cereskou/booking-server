package api

import (
	"ditto/booking/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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
	logon := logonFromToken(c)

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//input check
	if _, ok := input["name"]; !ok {
		return c.JSON(http.StatusBadRequest, BadRequest(errors.New("Tenant name is required")))
	}

	tx := s.DB().Begin()
	//create tenant and return
	tenant, err := s.DB().CreateTenant(tx, logon, input)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, InternalServerError(err))
	}

	//add logon user to tenant/ need re-login
	users := make([]*models.TenantUsers, 0)
	users = append(users, &models.TenantUsers{
		TenantID:   tenant.ID,
		UserID:     logon.ID,
		Right:      1,
		UpdateUser: logon.ID,
	})
	err = s.DB().AddUserToTenant(tx, logon, users)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, InternalServerError(err))
	}

	//clear cache
	s.CacheDel(logon.Email)

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
	logon := logonFromToken(c)

	sid := c.Param("id")
	tid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}
	tx := s.DB().Begin()

	err = s.DB().ChangeUserTenant(tx, logon, tid)
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
	logon := logonFromToken(c)

	tenants, err := s.DB().GetUserTenants(nil, logon, logon.ID)
	if err != nil {
		return err
	}
	resp := Response{
		Code: http.StatusOK,
		Data: tenants,
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
	logon := logonFromToken(c)

	users, err := s.DB().GetTenantUsesr(nil, logon, logon.Tenant)
	if err != nil {
		return err
	}

	resp := Response{
		Code: http.StatusOK,
		Data: users,
	}

	return c.JSON(http.StatusOK, resp)
}
