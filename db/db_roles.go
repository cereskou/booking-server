package db

import (
	"ditto/booking/cx"
	"ditto/booking/models"

	"gorm.io/gorm"
)

//GetRole -
func (d *Database) GetRole(db *gorm.DB, id int64, name string) (*models.Role, error) {
	db = d.ValidDB(db)

	ei := (id == 0)
	en := (name == "")

	var record models.Role
	result := db.Table(record.TableName()).
		Where("(id=? or ?) and (name=? or ?)", id, ei, name, en).
		First(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &record, nil
}

//GetRoles -
func (d *Database) GetRoles(db *gorm.DB) ([]*models.Role, error) {
	db = d.ValidDB(db)

	result := make([]*models.Role, 0)
	rows, err := db.Table("roles").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.Role
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

//DeleteRole -
func (d *Database) DeleteRole(db *gorm.DB, logon *cx.Payload, rid int64, name string) error {
	db = d.ValidDB(db)

	sql := "delete from roles where (id=? or ?) and (name=? or ?)"
	ei := (rid == 0)
	en := (name == "")
	err := db.Exec(sql, rid, ei, name, en).Error
	if err != nil {
		return err
	}
	return nil
}

//DeleteRoles -
func (d *Database) DeleteRoles(db *gorm.DB, logon *cx.Payload) error {
	db = d.ValidDB(db)

	sql := "delete from roles where id > 0"

	err := db.Exec(sql).Error
	if err != nil {
		return err
	}
	return nil
}
