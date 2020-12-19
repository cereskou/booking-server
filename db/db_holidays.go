package db

import (
	"ditto/booking/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//HolidaysInsert -
func (d *Database) HolidaysInsert(db *gorm.DB, recs []*models.Holiday) error {
	d.Lock()
	defer d.Unlock()
	db = d.ValidDB(db)

	values := make([]string, 0)
	for _, rec := range recs {
		val := fmt.Sprintf("(%q,%q,%v,%v)", rec.Ymd.Format("2006/1/2"), rec.Name, rec.Class, rec.UpdateUser)

		values = append(values, val)
	}
	sql := "insert into holidays(ymd,name,class,update_user) values " + strings.Join(values, ",") + " on duplicate key update name=values(name),,update_user=values(update_user)"
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

// HolidaysSelect -
func (d *Database) HolidaysSelect(db *gorm.DB, year string) ([]*models.Holiday, error) {
	db = d.ValidDB(db)

	result := make([]*models.Holiday, 0)
	rows, err := db.Table("holidays").Where("year(ymd)=?", year).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.Holiday
		if err := db.ScanRows(rows, &record); err != nil {
			break
		}

		result = append(result, &record)
	}
	if len(result) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return result, nil
}
