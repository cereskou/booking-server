package api

import (
	casbinmw "ditto/booking/casbin"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//RegisterRoutes -
func (s *Service) RegisterRoutes(e *echo.Echo, prefix string) {
	config := middleware.JWTConfig{
		SigningKey:    s._rsa.GetPublicKey(),
		SigningMethod: "RS512",
		TokenLookup:   "header:Authorization",
	}
	g := e.Group(prefix)

	g.POST("/login", s.Login)

	//refresh token
	g.POST("/refresh", s.Token)

	//Holiday
	g.GET("/holidays/:year", s.ListHolidays)
	h := g.Group("/holidays")
	h.Use(middleware.JWTWithConfig(config))
	h.Use(casbinmw.Middleware(s._enforcer))
	h.POST("", s.UpdateHolidays)

	u := g.Group("/member")
	u.Use(middleware.JWTWithConfig(config))
	u.Use(casbinmw.Middleware(s._enforcer))
	//User
	u.GET("/detail", s.User)

	//
	a := g.Group("/admin")
	a.Use(middleware.JWTWithConfig(config))
	a.Use(casbinmw.Middleware(s._enforcer))
	a.GET("/member/:email", s.Member)
}
