package api

import (
	"ditto/booking/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetUser - ログイン中ユーザー情報
// @Summary ログイン中ユーザー（自分）情報を取得します
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /user [get]
func (s *Service) GetUser(c echo.Context) error {
	logon := logonFromToken(c)

	user, err := s.DB().GetUser(logon.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp := Response{
				Code:  http.StatusNotFound,
				Error: err.Error(),
			}
			return c.JSON(http.StatusNotFound, resp)
		}
		return err
	}

	resp := Response{
		Code: http.StatusOK,
		Data: user,
	}

	return c.JSON(http.StatusOK, resp)

}

// GetAccount - ログイン中ユーザーのログイン情報
// @Summary ログイン中ユーザー（自分）ログイン情報を取得します
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /user/account [get]
func (s *Service) GetAccount(c echo.Context) error {
	logon := logonFromToken(c)

	account, err := s.DB().GetAccount(logon.Email)
	if err != nil {
		return echo.ErrBadRequest
	}
	account.PasswordHash = ""

	resp := Response{
		Code: http.StatusOK,
		Data: account,
	}

	return c.JSON(http.StatusOK, resp)

}

// UpdateUser - ユーザー情報を更新します
// @Summary ユーザー（自分）情報を更新します
// @Tags User
// @Accept json
// @Produce json
// @Param data body Empty true "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /user [put]
func (s *Service) UpdateUser(c echo.Context) error {
	logon := logonFromToken(c)

	input := make(map[string]interface{})
	//decode
	if err := c.Bind(&input); err != nil {
		return err
	}

	//update
	err := s.DB().UpdateUser(logon.ID, logon.Email, input)
	if err != nil {
		return err
	}

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}

// UpdatePassword - ユーザーパスワードを更新します
// @Summary ユーザー（自分）パスワードを更新します
// @Tags User
// @Accept json
// @Produce json
// @Param data body Password false "data"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} HTTPError
// @Security ApiKeyAuth
// @Router /user/password [put]
func (s *Service) UpdatePassword(c echo.Context) error {
	logon := logonFromToken(c)

	pwd := Password{}
	//decode
	if err := c.Bind(&pwd); err != nil {
		return err
	}

	//generate password hash code
	newHash := utils.HashPassword(pwd.NewPassword)

	//update password
	err := s.DB().UpdatePassword(logon.Email, newHash, pwd.UpdateDate)
	if err != nil {
		return err
	}

	//clear cache
	s.CacheDel(logon.Email)

	resp := Response{
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, resp)
}
