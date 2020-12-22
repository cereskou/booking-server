package api

import (
	"ditto/booking/models"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetDict - 辞書情報取得
// @Summary 辞書情報を取得します
// @Tags Dict
// @Accept json
// @Produce json
// @Param id path int true "辞書ID"
// @Param code path int false "コード"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /dict/{id}/{code} [get]
func (s *Service) GetDict(c echo.Context) error {
	//id
	dictid, err := paramInt(c, "id", "Dict id is required")
	if err != nil {
		return err
	}
	code, err := paramInt(c, "code", "Code is required")
	if err != nil {
		return err
	}

	resp := Response{
		Code: http.StatusOK,
	}
	if dictid == 0 {
		dicts, err := s.DB().GetAllDicts(nil)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(err)
			}
			return InternalServerError(err)
		}
		resp.Data = dicts
	} else {
		if code >= 0 {
			dict, err := s.DB().GetDict(nil, dictid, code)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return NotFound(err)
				}
				return InternalServerError(err)
			}
			resp.Data = dict
		} else {
			dicts, err := s.DB().GetDicts(nil, dictid)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return NotFound(err)
				}
				return InternalServerError(err)
			}
			resp.Data = dicts
		}
	}
	return c.JSON(http.StatusOK, resp)

}

// AddDict - 辞書情報作成
// @Summary 辞書情報を新規作成します
// @Tags Dict
// @Accept json
// @Produce json
// @Param data body Dict false "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /dict [post]
func (s *Service) AddDict(c echo.Context) error {
	logon := s.logonFromToken(c)

	data := Dict{}
	//decode
	if err := c.Bind(&data); err != nil {
		return err
	}

	if data.Code <= 0 {
		return BadRequest(errors.New("Code is 1-based integer"))
	}

	tx := s.DB().Begin()

	//update password
	rec := models.Dict{
		DictID:     int(data.DictID),
		Code:       int(data.Code),
		Kvalue:     data.Value,
		Remark:     data.Remark,
		Status:     data.Status,
		UpdateUser: logon.ID,
	}
	err := s.DB().AddDict(tx, logon, &rec)
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

// AddDicts - 辞書情報作成（複数）
// @Summary 辞書情報を新規作成します（複数）
// @Tags Dict
// @Accept json
// @Produce json
// @Param data body Dicts false "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /dict/array [post]
func (s *Service) AddDicts(c echo.Context) error {
	logon := s.logonFromToken(c)

	data := Dicts{}
	//decode
	if err := c.Bind(&data); err != nil {
		return err
	}

	values := make([]*models.Dict, 0)
	for _, v := range data.List {
		values = append(values, &models.Dict{
			TenantID:   logon.Tenant,
			DictID:     int(v.DictID),
			Code:       int(v.Code),
			Kvalue:     v.Value,
			Remark:     v.Remark,
			Status:     v.Status,
			UpdateUser: logon.ID,
		})
	}

	tx := s.DB().Begin()

	err := s.DB().AddDicts(tx, values)
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

// EnableDict - 辞書情報の有効無効
// @Summary 辞書情報を有効・無効します
// @Tags Dict
// @Accept json
// @Produce json
// @Param dictid path int true "辞書番号"
// @Param code query int false "コード"
// @Param status query int true "1:有効・0:無効"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /dict/{id}/enabled [put]
func (s *Service) EnableDict(c echo.Context) error {
	logon := s.logonFromToken(c)

	dictid, err := paramInt(c, "id", "Dict id is required")
	if err != nil {
		return err
	}
	code, err := queryParamInt(c, "code", "Code is required")
	if err != nil {
		return err
	}
	status, err := queryParamInt(c, "status", "Status is required")
	if err != nil {
		return err
	}

	tx := s.DB().Begin()

	if code > 0 {
		err = s.DB().EnableDict(tx, logon, dictid, code, int(status))
		if err != nil {
			tx.Rollback()
			return InternalServerError(err)
		}
	} else {
		err = s.DB().EnableDicts(tx, logon, dictid, int(status))
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

// UpdateDict - 辞書情報更新
// @Summary 辞書情報を更新します
// @Tags Dict
// @Accept json
// @Produce json
// @Param id path int true "辞書ID"
// @Param data body Dict true "データ"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /dict/{id} [put]
func (s *Service) UpdateDict(c echo.Context) error {
	logon := s.logonFromToken(c)

	dictid, err := paramInt(c, "id", "Dict id is required")
	if err != nil {
		return BadRequest(err)
	}

	data := Dict{}
	//decode
	if err := c.Bind(&data); err != nil {
		return err
	}

	if data.Code <= 0 {
		return BadRequest(errors.New("Code is 1-based integer"))
	}

	tx := s.DB().Begin()

	//update password
	rec := models.Dict{
		TenantID:   logon.Tenant,
		DictID:     int(dictid),
		Code:       int(data.Code),
		Kvalue:     data.Value,
		Remark:     data.Remark,
		Status:     data.Status,
		UpdateUser: logon.ID,
	}
	err = s.DB().UpdateDict(tx, logon, &rec)
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

// UpdateDicts - 辞書情報更新（複数）
// @Summary 辞書情報を更新します（複数）
// @Tags Dict
// @Accept json
// @Produce json
// @Param data body Dicts true "データ"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /dict/array [put]
func (s *Service) UpdateDicts(c echo.Context) error {
	return s.AddDicts(c)
}

// DeleteDict - 辞書情報削除
// @Summary 辞書情報を削除します（複数）
// @Tags Dict
// @Accept json
// @Produce json
// @Param id path int true "辞書番号"
// @Param code path int false "コード(0:全部)"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /dict/{id}/{code} [delete]
func (s *Service) DeleteDict(c echo.Context) error {
	logon := s.logonFromToken(c)

	//id
	dictid, err := paramInt(c, "id", "Dict id is required")
	if err != nil {
		return err
	}
	code, err := paramInt(c, "code", "Code is required")
	if err != nil {
		return err
	}

	tx := s.DB().Begin()

	if code >= 0 {
		err = s.DB().DeleteDict(tx, logon, dictid, code)
		if err != nil {
			tx.Rollback()
			return InternalServerError(err)
		}
	} else {
		err = s.DB().DeleteDicts(tx, logon, dictid)
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
