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
	//user confirm
	g.GET("/user/confirm", s.ConfirmEmail)

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
	//Logout
	u.GET("/logout", s.Logout)
	//get tenants
	u.GET("/tenants", s.GetTenants)
	u.PUT("/tenants/:id", s.ChangeUserTenant)

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
	a.GET("/user/:id", s.AdminGetUser)
	a.GET("/user/:id/account", s.AdminGetAccount)
	a.PUT("/user/:id", s.AdminUpdateUser)
	a.POST("/user", s.AdminCreateAccount)
	a.DELETE("/user/:id", s.AdminDeleteAcount)

	a.POST("/tenant", s.AdminCreateTenant)
	a.GET("/tenant", s.AdminGetTenant)
	a.PUT("/tenant/:id", s.AdminUpdateTenant)
	a.DELETE("/tenant/:id", s.AdminDeleteTenant)

	//Dict
	d := g.Group("/dict")
	d.Use(middleware.JWTWithConfig(config))
	d.Use(casbinmw.Middleware(s._enforcer))
	d.GET("", s.GetDict)
	d.POST("", s.AddDict)
	d.POST("/array", s.AddDicts)
	d.PUT("/:dictid/enabled", s.EnableDict)
	d.DELETE("", s.DeleteDict)

	//Tenant
	t := g.Group("/tenant")
	t.Use(middleware.JWTWithConfig(config))
	t.Use(casbinmw.Middleware(s._enforcer))
	t.GET("/users", s.TenantListUser)
	t.POST("", s.CreateTenant)
}
