package api

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

//paramInt - path
func paramInt(c echo.Context, name string, msg string) (int64, error) {
	param := c.Param(name)
	if param == "" {
		return 0, BadRequest(errors.New(msg))
	}
	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, BadRequest(err)
	}

	return value, nil
}

//queryParamInt - query
func queryParamInt(c echo.Context, name string, msg string) (int64, error) {
	param := c.QueryParam(name)
	if param == "" {
		return 0, BadRequest(errors.New(msg))
	}
	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, BadRequest(err)
	}

	return value, nil
}

func queryParam(c echo.Context, name string, msg string) (string, error) {
	param := c.QueryParam(name)
	if param == "" {
		return "", BadRequest(errors.New(msg))
	}

	return param, nil
}
