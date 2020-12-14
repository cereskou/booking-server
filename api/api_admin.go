package api

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

// Member - ユーザー情報
// @Summary ユーザー情報を取得します(admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Param email path string true "email"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/member/{email} [get]
func (s *Service) Member(c echo.Context) error {
	email := c.Param("email")
	email, _ = url.QueryUnescape(email)

	user, err := s.DB().GetUser(email)
	if err != nil {
		resp := Response{
			Code:  http.StatusNotFound,
			Error: err.Error(),
		}

		return c.JSON(http.StatusNotFound, resp)
	}

	resp := Response{
		Code: http.StatusOK,
		Data: user,
	}

	return c.JSON(http.StatusOK, resp)

}
