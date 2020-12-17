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

	//update password
	rec := models.Dict{
		DictID:     int(data.DictID),
		Code:       int(data.Code),
		Kvalue:     data.Value,
		Remark:     data.Remark,
		Status:     data.Status,
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
		values = append(values, &models.Dict{
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
		return err
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
// @Router /dict/{dictid}/enabled [put]
func (s *Service) EnableDict(c echo.Context) error {
	logon := logonFromToken(c)

	dictid, err := strconv.ParseInt(c.Param("dictid"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}
	code, err := strconv.ParseInt(c.QueryParam("code"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}
	status, err := strconv.ParseInt(c.QueryParam("status"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}

	tx := s.DB().Begin()

	if code >= 0 {
		err = s.DB().EnableDict(tx, logon.ID, dictid, code, int(status))
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, InternalServerError(err))
		}
	} else {
		err = s.DB().EnableDicts(tx, logon.ID, dictid, int(status))
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, InternalServerError(err))
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
// @Param data query Dict true "データ"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /dict/{dictid} [put]
func (s *Service) UpdateDict(c echo.Context) error {
	logon := logonFromToken(c)

	dictid, err := strconv.ParseInt(c.Param("dictid"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequest(err))
	}

	data := Dict{}
	//decode
	if err := c.Bind(&data); err != nil {
		return err
	}

	tx := s.DB().Begin()

	//update password
	rec := models.Dict{
		DictID:     int(dictid),
		Code:       int(data.Code),
		Kvalue:     data.Value,
		Remark:     data.Remark,
		Status:     data.Status,
		UpdateUser: logon.ID,
	}
	err = s.DB().UpdateDict(tx, logon.ID, &rec)
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

// UpdateDicts - 辞書情報更新（複数）
// @Summary 辞書情報を更新します（複数）
// @Tags Dict
// @Accept json
// @Produce json
// @Param data query Dicts true "データ"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /dict/array [put]
func (s *Service) UpdateDicts(c echo.Context) error {
	return s.AddDicts(c)
}
