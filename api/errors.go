package api

import (
	"ditto/booking/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

//Error -
type Error struct {
	StatusCode  int    `json:"code"`
	ResponsedAt string `json:"responsed_at"` // RFC3339
	Message     string `json:"message"`
}

//CustomHTTPErrorHandler -
func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := ""
	// https://godoc.org/github.com/labstack/echo#HTTPError
	if ee, ok := err.(*echo.HTTPError); ok {
		code = ee.Code
		message = ee.Message.(string)
	}
	body := Error{
		StatusCode:  code,
		ResponsedAt: utils.NowJST().String(),
		Message:     message,
	}

	c.JSON(code, body)
}
