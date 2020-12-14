package api

import (
	"ditto/booking/security"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

//Unauthorized -
func Unauthorized() *Response {
	return &Response{
		Code: http.StatusNotFound,
		Data: "Invalid Username or Password",
	}
}

//emailFromToken -
func emailFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	secret := claims["name"].(string)
	payload := security.DecryptString(secret)

	bodys := strings.Split(payload, "|")
	if len(bodys) == 3 {
		return bodys[1]
	}

	return ""
}
