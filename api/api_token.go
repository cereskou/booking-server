package api

import (
	"ditto/booking/security"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Token - api to refresh tokens
// @Summary リフレッシュトークン
// @Tags Account
// @Accept json
// @Produce json
// @Param data body RefreshToken false "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Router /refresh [post]
func (s *Service) Token(c echo.Context) error {
	req := RefreshToken{}

	//decode
	if err := c.Bind(&req); err != nil {
		return err
	}

	//refresh token
	if req.GrantType != "refresh_token" {
		return echo.ErrUnauthorized
	}

	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return s._rsa.GetPublicKey(), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		secret := claims["sub"].(string)
		payload := security.DecryptString(secret)

		bodys := strings.Split(payload, "|")
		if len(bodys) == 3 {
			name := bodys[0]
			email := bodys[1]
			role := bodys[2]

			//create a token
			token, err := s.generateToken(name, email, role)
			if err != nil {
				return err
			}
			resp := Response{
				Code: http.StatusOK,
				Data: token,
			}

			return c.JSON(http.StatusOK, resp)
		}

		return echo.ErrUnauthorized
	}

	return err
}
