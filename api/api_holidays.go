package api

import (
	"bufio"
	"bytes"
	"ditto/booking/config"
	"ditto/booking/logger"
	"ditto/booking/models"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// UpdateHolidays - 国民の祝日・休日を更新します
// @Summary 国民の祝日・休日を更新します
// @Tags Holidays
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /holidays [POST]
func (s *Service) UpdateHolidays(c echo.Context) error {
	conf := config.Load()

	logger.Trace(conf.Holidays.URL)

	client := resty.New()
	resp, err := client.R().Get(conf.Holidays.URL)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return err
	}

	//Encode convert
	b := bytes.NewReader(resp.Body())
	r := transform.NewReader(b, japanese.ShiftJIS.NewDecoder())
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	//JST
	jst, _ := time.LoadLocation("Asia/Tokyo")

	//Transaction
	var eoc bool = false
	tx := s.DB().Begin()
	defer func() {
		if eoc {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var count int64 = 0
	recs := make([]*models.Holiday, 0)
	var i int = 0
	b = bytes.NewReader(body)
	rx := bufio.NewReader(b)
	for {
		line, _, err := rx.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			eoc = true
			return err
		}

		if conf.Holidays.Header && i == 0 {
			i++
			continue
		}
		i++

		parts := strings.Split(string(line), ",")
		if len(parts) != 2 {
			continue
		}
		//convert string to time
		t, err := time.ParseInLocation("2006/1/2", parts[0], jst)
		if err != nil {
			continue
		}
		recs = append(recs, &models.Holiday{
			Ymd:        t,
			Name:       parts[1],
			Class:      0,
			UpdateUser: 1,
		})
		count++

		if len(recs) > 100 {
			err := s.DB().HolidaysInsert(tx, recs)
			if err != nil {
				eoc = true
				return err
			}
			recs = recs[:0]
		}
	}
	//余り分
	if len(recs) > 0 {
		err := s.DB().HolidaysInsert(tx, recs)
		if err != nil {
			eoc = true
			return err
		}
		recs = recs[:0]
	}

	result := Response{
		Code: http.StatusOK,
		Data: count,
	}

	return c.JSON(http.StatusOK, result)
}

// ListHolidays - ユーザー情報
// @Summary ユーザー情報を取得します(admin)
// @Tags Holidays
// @Accept json
// @Produce json
// @Param year path string true "year"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Router /holidays/{year} [get]
func (s *Service) ListHolidays(c echo.Context) error {
	year := c.Param("year")

	key := fmt.Sprintf("HOLI_%v", year)
	var holidays []*models.Holiday
	//use redis cache
	err := s.CacheGet(key, &holidays)
	if err != nil {
		holidays, err = s.DB().HolidaysSelect(nil, year)
		if err != nil {
			resp := Response{
				Code:  http.StatusNotFound,
				Error: err.Error(),
			}

			return c.JSON(http.StatusNotFound, resp)
		}
		err = s.CacheSet(key, holidays)
		if err != nil {
		}
	}

	resp := Response{
		Code: http.StatusOK,
		Data: holidays,
	}

	return c.JSON(http.StatusOK, resp)

}
