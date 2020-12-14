package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// User - ログイン中ユーザー情報
// @Summary ログイン中ユーザー（自分）情報を取得します
// @Tags Member
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /member/detail [get]
func (s *Service) User(c echo.Context) error {
	email := emailFromToken(c)

	user, err := s.DB().GetUser(email)
	if err != nil {
		return echo.ErrBadRequest
	}

	resp := Response{
		Code: http.StatusOK,
		Data: user,
	}

	return c.JSON(http.StatusOK, resp)

}
