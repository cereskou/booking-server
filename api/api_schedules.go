package api

import (
	"ditto/booking/models"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreateSchedule - スケジュールを作成
// @Summary スケジュールを作成します
// @Tags Schedule
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /schedule [post]
func (s *Service) CreateSchedule(c echo.Context) error {
	logon := s.logonFromToken(c)

	//input := make(map[string]interface{})
	data := models.Schedule{}
	//decode
	if err := c.Bind(&data); err != nil {
		return err
	}

	//input check
	tx := s.DB().Begin()
	//create
	facility, err := s.DB().CreateSchedule(tx, logon, &data)
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

// GetSchedules - スケジュール情報を取得（複数）
// @Summary スケジュール情報（複数）を取得します
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "メニューID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /menu/{id}/schedules [get]
func (s *Service) GetSchedules(c echo.Context) error {
	logon := s.logonFromToken(c)

	id, err := paramInt(c, "id", "Menu id is required")
	if err != nil {
		return err
	}

	result, err := s.DB().GetSchedules(nil, logon, id)
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

// GetSchedule - スケジュール情報を取得
// @Summary スケジュール情報を取得します
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "スケジュールID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /menu/schedule/{id} [get]
func (s *Service) GetSchedule(c echo.Context) error {
	logon := s.logonFromToken(c)

	//
	schedid, err := paramInt(c, "id", "Schedule id is required")
	if err != nil {
		return err
	}

	result, err := s.DB().GetSchedule(nil, logon, schedid)
	if err != nil {
		return InternalServerError(err)
	}

	resp := Response{
		Code: http.StatusOK,
		Data: result,
	}

	return c.JSON(http.StatusOK, resp)
}

// DeleteSchedule - スケジュールを削除します
// @Summary スケジュールを削除します
// @Tags Schedule
// @Accept json
// @Produce json
// @Param id path int true "スケジュールID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /schedule/{id} [delete]
func (s *Service) DeleteSchedule(c echo.Context) error {
	logon := s.logonFromToken(c)

	id, err := paramInt(c, "id", "Schedule id is required")
	if err != nil {
		return err
	}

	tx := s.DB().Begin()

	err = s.DB().DeleteSchedule(tx, logon, id)
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

// EnabledSchedule - スケジュールの利用可否（有効・無効）
// @Summary スケジュールの利用可否（有効・無効）
// @Tags Schedule
// @Accept json
// @Produce json
// @Param id path int true "スケジュールID"
// @Param status path int true "1:有効・0:無効"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /schedule/{id}/{status} [put]
func (s *Service) EnabledSchedule(c echo.Context) error {
	logon := s.logonFromToken(c)

	id, err := paramInt(c, "id", "Schedule id is required")
	if err != nil {
		return err
	}
	status, err := paramInt(c, "status", "Status is required")
	if err != nil {
		return err
	}

	tx := s.DB().Begin()

	err = s.DB().EnabledSchedule(tx, logon, id, int(status))
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
