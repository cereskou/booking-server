package api

import (
	"ditto/booking/cx"
	"ditto/booking/security"
	"ditto/booking/utils"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

//Unauthorized -
func Unauthorized() *Response {
	return &Response{
		Code: http.StatusUnauthorized,
		Data: "Invalid Username or Password",
	}
}

//BadRequest -
func BadRequest(err error) *Response {
	return &Response{
		Code:  http.StatusBadRequest,
		Error: err.Error(),
	}
}

//NotFound -
func NotFound(err error) *Response {
	return &Response{
		Code:  http.StatusNotFound,
		Error: err.Error(),
	}
}

//InternalServerError -
func InternalServerError(err error) *Response {
	return &Response{
		Code:  http.StatusInternalServerError,
		Error: err.Error(),
	}
}

//NewResponse -
func NewResponse(code int, msg string) *Response {
	return &Response{
		Code: code,
		Data: msg,
	}
}

//logonFromToken - get logon user from token
func logonFromToken(c echo.Context) *cx.Payload {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	secret := claims["uuid"].(string)
	payload := security.DecryptString(secret)

	var d cx.Payload
	err := utils.JSON.NewDecoder(strings.NewReader(payload)).Decode(&d)
	if err != nil {
		return nil
	}
	return &d
}
