package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreateMenu - メニューを作成
// @Summary メニューを作成します
// @Tags Menu
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /menu [post]
func (s *Service) CreateMenu(c echo.Context) error {
	logon := s.logonFromToken(c)

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//input check
	//name
	if _, ok := input["name"]; !ok {
		return BadRequest(errors.New("Menu name is required"))
	}

	tx := s.DB().Begin()
	//create
	facility, err := s.DB().CreateMenu(tx, logon, logon.Tenant, input)
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

// GetMenus - メニュー情報を取得（複数）
// @Summary メニュー情報を取得（複数）を取得します
// @Tags Tenant
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /tenant/menus [get]
func (s *Service) GetMenus(c echo.Context) error {
	logon := s.logonFromToken(c)

	result, err := s.DB().GetMenus(nil, logon, logon.Tenant)
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

// GetMenu - メニュー情報を取得
// @Summary メニュー情報を取得します
// @Tags Tenant
// @Accept json
// @Produce json
// @Param id path int true "メニューID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /tenant/menu/{id} [get]
func (s *Service) GetMenu(c echo.Context) error {
	logon := s.logonFromToken(c)

	if logon.Tenant == 0 {
		return NoContent(errors.New("No tenant"))
	}

	//class id
	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Menu id is required"))
	}
	cid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	result, err := s.DB().GetMenu(nil, logon, cid)
	if err != nil {
		return InternalServerError(err)
	}

	resp := Response{
		Code: http.StatusOK,
		Data: result,
	}

	return c.JSON(http.StatusOK, resp)
}

// UpdateMenu - メニュー情報を更新します
// @Summary メニュー情報を更新します
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "メニューID"
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /menu/{id} [put]
func (s *Service) UpdateMenu(c echo.Context) error {
	logon := s.logonFromToken(c)

	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Menu id is required"))
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

	err = s.DB().UpdateMenu(tx, logon, id, input)
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

// DeleteMenu - メニューを削除します
// @Summary メニューを削除します
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "メニューID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /menu/{id} [delete]
func (s *Service) DeleteMenu(c echo.Context) error {
	logon := s.logonFromToken(c)

	sdid := c.Param("id")
	if sdid == "" {
		return BadRequest(errors.New("Menu id is required"))
	}
	id, err := strconv.ParseInt(sdid, 10, 64)
	if err != nil {
		return BadRequest(err)
	}

	tx := s.DB().Begin()

	err = s.DB().DeleteMenu(tx, logon, id)
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

// EnabledMenu - メニューの利用可否（有効・無効）
// @Summary メニューの利用可否（有効・無効）
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "メニューID"
// @Param status path int true "1:有効・0:無効"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /menu/{id}/{status} [put]
func (s *Service) EnabledMenu(c echo.Context) error {
	logon := s.logonFromToken(c)

	sid := c.Param("id")
	if sid == "" {
		return BadRequest(errors.New("Menu id is required"))
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

	err = s.DB().EnabledMenu(tx, logon, id, int(status))
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
