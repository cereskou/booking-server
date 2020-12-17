package db

import "gorm.io/gorm"

//MyCallback -
type MyCallback func(map[string]interface{}) error

//ValidDB -
func (d *Database) ValidDB(db *gorm.DB) *gorm.DB {
	if db == nil {
		return d.DB()
	}

	return db
}
