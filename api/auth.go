package api

import (
	"ditto/booking/cx"
	"ditto/booking/utils"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

//logonFromToken - get logon user from token
func (s *Service) logonFromToken(c echo.Context) *cx.Payload {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	secret := claims["uuid"].(string)
	//payload := security.DecryptString(secret)
	payload := secret

	var d cx.Payload
	err := utils.JSON.NewDecoder(strings.NewReader(payload)).Decode(&d)
	if err != nil {
		return nil
	}

	//set tenant id
	err = s.reacquireTenant(&d)
	if err != nil {
	}

	return &d
}

func (s *Service) reacquireTenant(logon *cx.Payload) error {
	var id int64

	key := "ACCN_TENANT_" + logon.Email
	var t struct {
		TenantID int64 `json:"tenant_id"` //テナントID
		UserID   int64 `json:"user_id"`   //ユーザーID
		Right    int   `json:"right"`     //権限
	}
	err := s.CacheGet(key, &t)
	if err != nil {
		t, err := s.DB().GetUserTenant(nil, logon.ID)
		if err != nil {
			return err
		}
		id = t.TenantID
		err = s.CacheSet(key, t)
		if err != nil {

		}
	} else {
		id = t.TenantID
	}
	//set tenant id
	logon.Tenant = id

	return nil
}
