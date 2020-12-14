package db

import (
	"sync"
	"time"

	"ditto/booking/config"
	"ditto/booking/logger"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//Database - sqlite3 database
type Database struct {
	sync.Mutex
	Error error
	Name  string
	_db   *gorm.DB
}

//NewDatabase - cereate a new database
func NewDatabase() *Database {
	dbase := &Database{
		Error: nil,
	}
	conf := config.Load()

	logger.Trace("DNS: ", conf.Db.DNS)
	db, err := gorm.Open(mysql.Open(conf.Db.DNS), &gorm.Config{})
	if err != nil {
		dbase.Error = err

		return dbase
	}
	//trace
	logger.Trace("MaxOpenConns ", conf.Db.MaxOpenConns)
	logger.Trace("MaxIdleConns ", conf.Db.MaxIdleConns)
	logger.Trace("ConnMaxLifetime ", time.Second*time.Duration(conf.Db.ConnMaxLifetime))
	logger.Trace("Debug ", conf.Db.Debug)

	sqldb, _ := db.DB()
	// SetMaxOpenConnsは接続済みのデータベースコネクションの最大数を設定します
	// SetMaxIdleConns()はSetMaxOpenConns()以上に
	sqldb.SetMaxOpenConns(conf.Db.MaxOpenConns)
	// SetMaxIdleConnsはアイドル状態のコネクションプール内の最大数を設定します
	sqldb.SetMaxIdleConns(conf.Db.MaxIdleConns)
	// SetConnMaxLifetime() は最大接続数×1秒 程度に
	// SetConnMaxLifetimeは再利用され得る最長時間を設定します
	sqldb.SetConnMaxLifetime(time.Second * time.Duration(conf.Db.ConnMaxLifetime))

	//set
	dbase._db = db
	dbase.Error = nil
	if conf.Db.Debug {
		dbase._db = db.Debug()
	}

	return dbase
}

//Close -
func (d *Database) Close() error {
	logger.Debug("database close")
	db, err := d._db.DB()
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}

//DB -
func (d *Database) DB() *gorm.DB {
	return d._db
}

//Begin -
func (d *Database) Begin() *gorm.DB {
	return d._db.Begin()
}

//Rollback -
func (d *Database) Rollback() *gorm.DB {
	return d._db.Rollback()
}

//Commit -
func (d *Database) Commit() *gorm.DB {
	return d._db.Commit()
}
