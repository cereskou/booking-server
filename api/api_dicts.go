package api

import (
	"ditto/booking/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetDict - 辞書情報取得
// @Summary 辞書情報を取得します
// @Tags Dict
// @Accept json
// @Produce json
// @Param dictid query int true "辞書ID"
// @Param code query int false "コード"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /dict [get]
func (s *Service) GetDict(c echo.Context) error {
	dictid, err := strconv.ParseInt(c.QueryParam("dictid"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}
	code, err := strconv.ParseInt(c.QueryParam("code"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}

	resp := Response{
		Code: http.StatusOK,
	}

	if dictid == 0 {
		dicts, err := s.DB().GetAllDicts(nil)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.JSON(http.StatusNotFound, NotFound(err))
			}
			return c.JSON(http.StatusInternalServerError, InternalServerError(err))
		}
		resp.Data = dicts
	} else {
		if code >= 0 {
			dict, err := s.DB().GetDict(nil, dictid, code)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return c.JSON(http.StatusNotFound, NotFound(err))
				}
				return c.JSON(http.StatusInternalServerError, InternalServerError(err))
			}
			resp.Data = dict
		} else {
			dicts, err := s.DB().GetDicts(nil, dictid)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return c.JSON(http.StatusNotFound, NotFound(err))
				}
				return c.JSON(http.StatusInternalServerError, InternalServerError(err))
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
	logon := logonFromToken(c)

	data := Dict{}
	//decode
	if err := c.Bind(&data); err != nil {
		return err
	}

	tx := s.DB().Begin()

	status := 1
	if !data.Status {
		status = 0
	}
	//update password
	rec := models.Dict{
		DictID:     int(data.DictID),
		Code:       int(data.Code),
		Kvalue:     data.Value,
		Remark:     data.Remark,
		Status:     status,
		UpdateUser: logon.ID,
	}
	err := s.DB().AddDict(tx, &rec)
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
	logon := logonFromToken(c)

	data := Dicts{}
	//decode
	if err := c.Bind(&data); err != nil {
		return err
	}

	values := make([]*models.Dict, 0)
	for _, v := range data.Dict {
		status := 0
		if v.Status {
			status = 1
		}
		values = append(values, &models.Dict{
			DictID:     int(v.DictID),
			Code:       int(v.Code),
			Kvalue:     v.Value,
			Remark:     v.Remark,
			Status:     status,
			UpdateUser: logon.ID,
		})
	}

	tx := s.DB().Begin()

	err := s.DB().AddDicts(tx, values)
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
