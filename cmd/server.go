package cmd

import (
	"ditto/booking/config"
	"ditto/booking/db"
	"ditto/booking/logger"
	"ditto/booking/services"
	"fmt"
	"net/http"
	"time"

	_ "ditto/booking/docs"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tylerb/graceful"

	echoSwagger "github.com/swaggo/echo-swagger"
)

//RunServer -
func RunServer(db *db.Database) error {
	logger.Info("Run Server")

	//get config
	conf := config.Load()

	var client *redis.Client
	//Setup cache
	if conf.Cache.Enable {
		client = redis.NewClient(&redis.Options{
			Addr:     conf.Cache.Address,
			Password: conf.Cache.Password,
			DB:       0, //use default DB
		})
		//check the server is alive
		if err := client.Ping().Err(); err != nil {
			logger.Error(err)
		}
	}

	//InitService
	if err := services.InitService(db, client); err != nil {
		return err
	}

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		if client != nil {
			val, err := client.Get("BOOKING").Result()
			if err != nil {
				logger.Error(err)

				return c.String(http.StatusOK, "Not found")
			}
			return c.String(http.StatusOK, val)
		}

		return c.String(http.StatusOK, "No redis")
	})

	//Router
	services.RegisterRoutes(e)

	//Server address
	e.Server.Addr = fmt.Sprintf("%v:%v", conf.Host, conf.Port)
	timeout := time.Duration(conf.Timeout) * time.Second
	logger.Info(e.Server.Addr)

	//Serve
	graceful.ListenAndServe(e.Server, timeout)

	logger.Info("Stop")

	return nil
}
