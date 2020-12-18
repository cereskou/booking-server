package mw

import (
	"ditto/booking/security"
	"ditto/booking/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//CasbinConfig -
type CasbinConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// Enforcer CasbinAuth main rule.
	// Required.
	Enforcer *casbin.Enforcer
}

var (
	// DefaultConfig is the default CasbinAuth middleware config.
	DefaultConfig = CasbinConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

// Middleware returns a CasbinAuth middleware.
//
// For valid credentials it calls the next handler.
// For missing or invalid credentials, it sends "401 - Unauthorized" response.
func Middleware(ce *casbin.Enforcer) echo.MiddlewareFunc {
	c := DefaultConfig
	c.Enforcer = ce
	return MiddlewareWithConfig(c)
}

// MiddlewareWithConfig returns a CasbinAuth middleware with config.
// See `Middleware()`.
func MiddlewareWithConfig(config CasbinConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if pass, err := config.CheckPermission(c); err == nil && pass {
				return next(c)
			} else if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return echo.ErrForbidden
		}
	}
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *CasbinConfig) CheckPermission(c echo.Context) (bool, error) {
	role := "anonymous"
	user := c.Get("user")
	if user != nil {
		role = roleFromToken(c)
	}
	method := c.Request().Method
	path := c.Request().URL.Path

	//複数対応
	r := strings.Split(role, ",")
	if len(r) > 1 {
		for _, r0 := range r {
			b, _ := a.Enforcer.Enforce(r0, path, method)
			if b {
				return b, nil
			}
		}
		return false, errors.New("Forbidden")
	}

	return a.Enforcer.Enforce(role, path, method)
}

//decode role from token
func roleFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	secret := claims["uuid"].(string)

	payload := security.DecryptString(secret)
	var d struct {
		ID       int64  `json:"id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Role     string `json:"role"`
		TenantID int64  `json:"tenantid"`
	}
	err := utils.JSON.NewDecoder(strings.NewReader(payload)).Decode(&d)
	if err != nil {
		return ""
	}
	return d.Role
}
