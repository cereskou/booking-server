package main

import (
	"ditto/booking/cmd"
	"ditto/booking/config"
	"ditto/booking/db"
	"ditto/booking/logger"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

//options - パラメータ定義
type options struct {
	Port int    `short:"p" long:"port" description:"server port"`
	Cmd  string `long:"run" description:"run command[migration, server]" default:"server"`
	Log  string `short:"l" long:"log" description:"Log level (trace, debug, warn, info, error)"`
}

// @title Booking
// @version 1.0
// @description 予約システム
// @termsOfService http://ditto.co.jp/terms/

// @contact.name API Support
// @contact.url https://www.ditto.co.jp/support
// @contact.email support@ditto.co.jp

// @license.name ditto license
// @license.url https://www.ditto.co.jp
// @host localhost:4000
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	//Get parameter
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(-1)
	}
	//Load config
	conf := config.Load()
	if conf.Error != nil {
		fmt.Println(conf.Error)
		os.Exit(-1)
	}
	conf.Version = version

	//log level
	if opts.Log != "" {
		conf.Log.Level = opts.Log
	}
	//set log level
	logger.SetLevel(conf.Log.Level)

	// initialize database
	logger.Debug("Initialize database")
	//new database instance
	db := db.NewDatabase()
	if db == nil || db.Error != nil {
		logger.Error(db.Error)
		os.Exit(-1)
	}
	defer db.Close()

	var err error
	switch strings.ToLower(opts.Cmd) {
	case "server":
		err = cmd.RunServer(db)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	case "migration":
		err = cmd.Migrate(db)
		if err != nil {
			logger.Error(err)
			os.Exit(-1)
		}
	}
}
