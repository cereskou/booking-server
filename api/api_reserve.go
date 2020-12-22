package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreateReservation - 予約を作成
// @Summary 予約を作成します
// @Tags Reservation
// @Accept json
// @Produce json
// @Param schedid path int true "スケジュールID"
// @Param userid path int true "ユーザーID"
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /reserve/{schedid}/{userid} [post]
func (s *Service) CreateReservation(c echo.Context) error {
	logon := s.logonFromToken(c)

	sid, err := paramInt(c, "schedid", "Schedule id is required")
	if err != nil {
		return err
	}
	uid, err := paramInt(c, "userid", "User id is required")
	if err != nil {
		return err
	}
	data := make(map[string]interface{})
	//decode
	if err := c.Bind(&data); err != nil {
		return err
	}

	//input check
	tx := s.DB().Begin()
	//create
	result, err := s.DB().CreateReservation(tx, logon, sid, uid, data)
	if err != nil {
		tx.Rollback()
		return InternalServerError(err)
	}
	tx.Commit()

	resp := Response{
		Code: http.StatusOK,
		Data: result,
	}

	return c.JSON(http.StatusOK, resp)
}

// GetUserReservations - 予約情報を取得（複数）
// @Summary 予約情報（複数）を取得します
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /user/reserve [get]
func (s *Service) GetUserReservations(c echo.Context) error {
	logon := s.logonFromToken(c)

	result, err := s.DB().GetReservations(nil, logon, 0, logon.ID)
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

// GetScheduleReservations - 予約情報を取得（複数）
// @Summary 予約情報（複数）を取得します
// @Tags Schedule
// @Accept json
// @Produce json
// @Param id path int true "スケジュールID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /schedule/{id}/reserve [get]
func (s *Service) GetScheduleReservations(c echo.Context) error {
	logon := s.logonFromToken(c)

	//id
	schedid, err := paramInt(c, "id", "Reservation id is required")
	if err != nil {
		return err
	}

	result, err := s.DB().GetReservations(nil, logon, schedid, 0)
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

// GetReservation - 予約情報を取得
// @Summary 予約情報を取得します
// @Tags Reservation
// @Accept json
// @Produce json
// @Param id path int true "予約ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /reserve/{id} [get]
func (s *Service) GetReservation(c echo.Context) error {
	logon := s.logonFromToken(c)

	//id
	id, err := paramInt(c, "id", "Reservation id is required")
	if err != nil {
		return err
	}

	result, err := s.DB().GetReservation(nil, logon, id)
	if err != nil {
		return InternalServerError(err)
	}

	resp := Response{
		Code: http.StatusOK,
		Data: result,
	}

	return c.JSON(http.StatusOK, resp)
}

// DeleteReservation - 予約を削除します
// @Summary 予約を削除します
// @Tags Reservation
// @Accept json
// @Produce json
// @Param id path int true "予約ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /reserve/{id} [delete]
func (s *Service) DeleteReservation(c echo.Context) error {
	logon := s.logonFromToken(c)

	id, err := paramInt(c, "id", "Reservation id is required")
	if err != nil {
		return err
	}

	tx := s.DB().Begin()

	err = s.DB().DeleteReservation(tx, logon, id)
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

// EnabledReservation - 予約の利用可否（有効・無効）
// @Summary 予約の利用可否（有効・無効）
// @Tags Reservation
// @Accept json
// @Produce json
// @Param id path int true "予約ID"
// @Param status path int true "1:有効・0:無効"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /reserve/{id}/{status} [put]
func (s *Service) EnabledReservation(c echo.Context) error {
	logon := s.logonFromToken(c)

	id, err := paramInt(c, "id", "Reservation id is required")
	if err != nil {
		return err
	}
	status, err := paramInt(c, "status", "Status is required")
	if err != nil {
		return err
	}

	tx := s.DB().Begin()

	err = s.DB().EnabledReservation(tx, logon, id, int(status))
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
