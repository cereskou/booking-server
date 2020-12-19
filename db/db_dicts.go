package db

import (
	"ditto/booking/cx"
	"ditto/booking/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//GetDict -
func (d *Database) GetDict(db *gorm.DB, dictid int64, code int64) (*models.Dict, error) {
	db = d.ValidDB(db)

	var record models.Dict
	result := db.Table(record.TableName()).
		Where("dict_id=? and code=? and status=1", dictid, code).
		First(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &record, nil
}

//GetDicts -
func (d *Database) GetDicts(db *gorm.DB, dictid int64) ([]*models.Dict, error) {
	db = d.ValidDB(db)

	result := make([]*models.Dict, 0)
	rows, err := db.Table("dicts").Where("dict_id=? and status=1", dictid).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.Dict
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

//GetAllDicts -
func (d *Database) GetAllDicts(db *gorm.DB) ([]*models.Dict, error) {
	db = d.ValidDB(db)

	result := make([]*models.Dict, 0)
	rows, err := db.Table("dicts").Where("status=1").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.Dict
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

//UpdateDict -
func (d *Database) UpdateDict(db *gorm.DB, logon *cx.Payload, data *models.Dict) error {
	db = d.ValidDB(db)

	sql := "update dicts set update_user=?,kvalue=? where dict=? and code=?"

	err := db.Exec(sql, logon.ID, data.Kvalue, data.DictID, data.Code).Error
	if err != nil {
		return err
	}
	return nil
}

//EnableDict -
func (d *Database) EnableDict(db *gorm.DB, logon *cx.Payload, dictid int64, code int64, status int) error {
	db = d.ValidDB(db)

	sql := "update dicts set update_user=?,status=? where dict=? and code=?"

	err := db.Exec(sql, logon.ID, status, dictid, code).Error
	if err != nil {
		return err
	}
	return nil
}

//EnableDicts -
func (d *Database) EnableDicts(db *gorm.DB, logon *cx.Payload, dictid int64, status int) error {
	db = d.ValidDB(db)

	sql := "update dicts set update_user=?,status=? where dict=?"

	err := db.Exec(sql, logon.ID, status, dictid).Error
	if err != nil {
		return err
	}
	return nil
}

//DeleteDict -
func (d *Database) DeleteDict(db *gorm.DB, logon *cx.Payload, dictid int64, code int64) error {
	db = d.ValidDB(db)

	sql := "delete from dicts where tenant_id=? and dict_id=? and code=?"

	err := db.Exec(sql, logon.Tenant, dictid, code).Error
	if err != nil {
		return err
	}
	return nil
}

//DeleteDicts -
func (d *Database) DeleteDicts(db *gorm.DB, logon *cx.Payload, dictid int64) error {
	db = d.ValidDB(db)

	sql := "delete from dicts where tenant_id=? and dict_id=?"

	err := db.Exec(sql, logon.Tenant, dictid).Error
	if err != nil {
		return err
	}
	return nil
}

//AddDict -
func (d *Database) AddDict(db *gorm.DB, logon *cx.Payload, rec *models.Dict) error {
	db = d.ValidDB(db)

	sql := "insert into dicts(dict_id,code,kvalue,remark,status,update_user) values (?,?,?,?,?,?)"
	result := db.Exec(sql, rec.DictID, rec.Code, rec.Kvalue, rec.Remark, rec.Status, rec.UpdateUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

//AddDicts -
func (d *Database) AddDicts(db *gorm.DB, list []*models.Dict) error {
	db = d.ValidDB(db)

	values := make([]string, 0)
	for _, rec := range list {
		val := fmt.Sprintf("(%d,%d,%q,%q,%d,%v)", rec.DictID, rec.Code, rec.Kvalue, rec.Remark, rec.Status, rec.UpdateUser)

		values = append(values, val)
	}
	sql := "insert into dicts(dict_id,code,kvalue,remark,status,update_user) values " + strings.Join(values, ",") +
		" on duplicate key update kvalue=values(kvalue),remark=values(remark),status=values(status),update_user=values(update_user)"
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}
