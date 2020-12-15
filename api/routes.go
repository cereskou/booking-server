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

	//anonymous
	//load
	g.POST("/user/login", s.Login)
	//refresh token
	g.POST("/user/refresh", s.RefreshToken)

	//User
	u := g.Group("/user")
	u.Use(middleware.JWTWithConfig(config))
	u.Use(casbinmw.Middleware(s._enforcer))
	//get user detail
	u.GET("", s.GetUser)
	//get login detail
	u.GET("/account", s.GetAccount)
	//update user
	u.PUT("", s.UpdateUser)
	//update user password
	u.PUT("/password", s.UpdatePassword)

	//Holiday
	g.GET("/holidays/:year", s.ListHolidays)
	h := g.Group("/holidays")
	h.Use(middleware.JWTWithConfig(config))
	h.Use(casbinmw.Middleware(s._enforcer))
	h.POST("", s.UpdateHolidays)

	//Admin
	a := g.Group("/admin")
	a.Use(middleware.JWTWithConfig(config))
	a.Use(casbinmw.Middleware(s._enforcer))
	a.GET("/user/:email", s.AdminGetUser)
	a.GET("/account/:email", s.AdminGetAccount)
	a.PUT("/user/:email", s.AdminUpdateUser)
}
