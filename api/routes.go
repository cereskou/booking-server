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
	//signup
	g.POST("/signup", s.Signup)
	//load
	g.POST("/login", s.Login)
	//refresh token
	g.POST("/refresh", s.RefreshToken)
	//user confirm
	g.GET("/confirm", s.ConfirmEmail)

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
	//get classes
	u.GET("/classes", s.GetClasses)
	//get reservations
	u.GET("/reserve", s.GetUserReservations)

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
	d.GET("/:id/:code", s.GetDict)
	d.POST("", s.AddDict)
	d.POST("/array", s.AddDicts)
	d.PUT("/:id/enabled", s.EnableDict)
	d.DELETE("/:id/:code", s.DeleteDict)
	d.PUT("/:id", s.UpdateDict)

	//Tenant -
	t := g.Group("/tenant")
	t.Use(middleware.JWTWithConfig(config))
	t.Use(casbinmw.Middleware(s._enforcer))
	t.GET("/users", s.TenantListUser)
	t.GET("/users/detail", s.TenantListUserWithDetail)
	t.POST("", s.CreateTenant)
	//create a user
	t.POST("/user", s.TenantCreateUser)
	//add/remove exist user to tenant
	t.PUT("/user", s.TenantDividedUser)
	t.DELETE("/user/:id", s.TenantDeleteUser)
	// tenant/facility
	t.GET("/facilities", s.GetFacilities)
	t.GET("/facility/:id", s.GetFacility)
	// tenant/menu
	t.GET("/menus", s.GetMenus)
	t.GET("/menu/:id", s.GetMenu)

	//Class -
	c := g.Group("/class")
	c.Use(middleware.JWTWithConfig(config))
	c.Use(casbinmw.Middleware(s._enforcer))
	c.GET("/users/:id", s.ClassListUser)
	c.GET("/users/:id/detail", s.ClassListUserWithDetail)
	c.POST("", s.CreateClass)
	c.POST("/user/:id", s.ClassCreateUser)
	c.PUT("/user/:id", s.ClassDividedUser)

	//facility -
	f := g.Group("/facility")
	f.Use(middleware.JWTWithConfig(config))
	f.Use(casbinmw.Middleware(s._enforcer))
	f.POST("", s.CreateFacility)
	f.PUT("/:id", s.UpdateFacility)
	f.DELETE("/:id", s.DeleteFacility)
	f.PUT("/:id/:status", s.EnabledFacility)

	//menu -
	m := g.Group("/menu")
	m.Use(middleware.JWTWithConfig(config))
	m.Use(casbinmw.Middleware(s._enforcer))
	m.POST("", s.CreateMenu)
	m.PUT("/:id", s.UpdateMenu)
	m.DELETE("/:id", s.DeleteMenu)
	m.PUT("/:id/:status", s.EnabledMenu)

	m.GET("/:id/schedules", s.GetSchedules)
	m.GET("/schedule/:id", s.GetSchedule)

	//schedule -
	l := g.Group("/schedule")
	l.Use(middleware.JWTWithConfig(config))
	l.Use(casbinmw.Middleware(s._enforcer))
	l.POST("", s.CreateSchedule)
	l.DELETE("/:id", s.DeleteSchedule)
	l.PUT("/:id/:status", s.EnabledSchedule)
	l.GET("/:id/reserve", s.GetScheduleReservations)

	//reserve -
	r := g.Group("/reserve")
	r.Use(middleware.JWTWithConfig(config))
	r.Use(casbinmw.Middleware(s._enforcer))
	r.POST("/:schedid/:userid", s.CreateReservation)
	r.GET("/:id", s.GetReservation)
	r.DELETE("/:id", s.DeleteReservation)
	r.PUT("/:id/:status", s.EnabledReservation)
}

//traceID -
func (s *Service) traceID(c echo.Context) string {
	return c.Request().Header.Get("x-trace")
}

//spanID -
func (s *Service) spanID(c echo.Context) string {
	return c.Request().Header.Get("x-span")
}
