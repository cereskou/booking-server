package services

import (
	"ditto/booking/api"
	"ditto/booking/config"
	"ditto/booking/db"
	"ditto/booking/logger"
	"ditto/booking/rsa"
	"reflect"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
)

var (
	//rsa
	_rsa *rsa.RSA

	//APIService ...
	APIService api.ServiceInterface
)

//InitService -
func InitService(db *db.Database, client *redis.Client) error {
	logger.Debug("Services InitService")
	conf := config.Load()

	logger.Trace("RSA Private: ", conf.Rsa.Private)
	logger.Trace("RSA Public : ", conf.Rsa.Public)

	var err error
	_rsa, err = rsa.NewRSA(conf.Rsa.Private, conf.Rsa.Public)
	if err != nil {
		return err
	}

	//New service
	if nil == reflect.TypeOf(APIService) {
		APIService, err = api.New(db, _rsa, client)
		if err != nil {
			return err
		}
	}

	return nil
}

//Close -
func Close() {
	logger.Debug("Services Close")
	APIService.Close()
}

//RegisterRoutes -
func RegisterRoutes(e *echo.Echo) {
	conf := config.Load()
	//API
	APIService.RegisterRoutes(e, conf.BaseURL)
}
