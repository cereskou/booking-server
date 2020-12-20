package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreateFacility - 施設を作成
// @Summary 施設を作成します
// @Tags Facility
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /facility [post]
func (s *Service) CreateFacility(c echo.Context) error {
	logon := s.logonFromToken(c)

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//input check
	//name
	if _, ok := input["name"]; !ok {
		return BadRequest(errors.New("Facility name is required"))
	}

	tx := s.DB().Begin()
	//create
	facility, err := s.DB().CreateFacility(tx, logon, logon.Tenant, input)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}
	tx.Commit()

	resp := Response{
		Code: http.StatusOK,
		Data: facility,
	}

	return c.JSON(http.StatusOK, resp)
}

// GetFacilities - 施設情報を取得（複数）
// @Summary 施設情報を取得（複数）を取得します
// @Tags Tenant
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /tenant/facilities [get]
func (s *Service) GetFacilities(c echo.Context) error {
	logon := s.logonFromToken(c)

	result, err := s.DB().GetFacilities(nil, logon, logon.Tenant)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return NotFound(err)
		}
		return InternalServerError(err)
	}
	if len(result) == 0 {
		return NotFound(errors.New("No content"))
	}

	resp := Response{
		Code: http.StatusOK,
		Data: result,
	}

	return c.JSON(http.StatusOK, resp)
}

// GetFacility - 施設情報を取得
// @Summary 施設情報を取得します
// @Tags Tenant
// @Accept json
// @Produce json
// @Param id path int true "施設ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /tenant/facility/{id} [get]
func (s *Service) GetFacility(c echo.Context) error {
	logon := s.logonFromToken(c)

	if logon.Tenant == 0 {
		return NoContent(errors.New("No tenant"))
	}

	//class id
	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Facility id is required"))
	}
	cid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	result, err := s.DB().GetFacility(nil, logon, cid)
	if err != nil {
		return InternalServerError(err)
	}

	resp := Response{
		Code: http.StatusOK,
		Data: result,
	}

	return c.JSON(http.StatusOK, resp)
}

// UpdateFacility - 施設情報を更新します
// @Summary 施設情報を更新します
// @Tags Facility
// @Accept json
// @Produce json
// @Param id path int true "施設ID"
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /facility/{id} [put]
func (s *Service) UpdateFacility(c echo.Context) error {
	logon := s.logonFromToken(c)

	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Facility id is required"))
	}
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	tx := s.DB().Begin()

	err = s.DB().UpdateFacility(tx, logon, id, input)
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

// DeleteFacility - 施設を削除します
// @Summary 施設を削除します
// @Tags Facility
// @Accept json
// @Produce json
// @Param id path int true "施設ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /facility/{id} [delete]
func (s *Service) DeleteFacility(c echo.Context) error {
	logon := s.logonFromToken(c)

	sdid := c.Param("id")
	if sdid == "" {
		return BadRequest(errors.New("Facility id is required"))
	}
	id, err := strconv.ParseInt(sdid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	tx := s.DB().Begin()

	err = s.DB().DeleteFacility(tx, logon, id)
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

// EnabledFacility - 施設の利用可否（有効・無効）
// @Summary 施設の利用可否（有効・無効）
// @Tags Facility
// @Accept json
// @Produce json
// @Param id path int true "施設ID"
// @Param status path int true "1:有効・0:無効"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /facility/{id}/{status} [put]
func (s *Service) EnabledFacility(c echo.Context) error {
	logon := s.logonFromToken(c)

	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Facility id is required"))
	}
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}
	ssta := c.Param("status")
	if ssta == "" {
		return BadRequest(errors.New("Status is required"))
	}
	status, err := strconv.ParseInt(ssta, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	tx := s.DB().Begin()

	err = s.DB().EnabledFacility(tx, logon, id, int(status))
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
