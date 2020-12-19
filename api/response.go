package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//Unauthorized - 401
func Unauthorized() *Response {
	return &Response{
		Code: http.StatusUnauthorized,
		Data: "Invalid Username or Password",
	}
}

//BadRequest - 400
func BadRequest(err error) error {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

//NotFound - 404
func NotFound(err error) error {
	return echo.NewHTTPError(http.StatusNotFound, err.Error())
}

//NoContent - 204
func NoContent(err error) error {
	return echo.NewHTTPError(http.StatusNoContent, err.Error())
}

//InternalServerError - 500
func InternalServerError(err error) error {
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}

//NewResponse -
func NewResponse(code int, msg string) *Response {
	return &Response{
		Code: code,
		Data: msg,
	}
}
