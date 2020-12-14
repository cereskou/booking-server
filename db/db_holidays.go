package db

import (
	"ditto/booking/models"
	"fmt"
	"strings"
)

//HolidaysInsert -
func (d *Database) HolidaysInsert(recs []*models.Holiday) error {
	d.Lock()
	defer d.Unlock()

	values := make([]string, 0)
	for _, rec := range recs {
		val := fmt.Sprintf("(%q,%q,%v,%v)", rec.Ymd.Format("2006/1/2"), rec.Name, rec.Class, rec.UpdateUser)

		values = append(values, val)
	}
	sql := "insert into holidays(ymd,name,class,update_user) values " + strings.Join(values, ",") + " on duplicate key update name = values(name)"
	err := d.DB().Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

// HolidaysSelect -
func (d *Database) HolidaysSelect(year string) ([]*models.Holiday, error) {
	result := make([]*models.Holiday, 0)
	rows, err := d.DB().Table("holidays").Where("year(ymd)=?", year).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.Holiday
		if err := d.DB().ScanRows(rows, &record); err != nil {
			break
		}

		result = append(result, &record)
	}

	return result, nil
}
