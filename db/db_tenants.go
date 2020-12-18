package db

import (
	"ditto/booking/cx"
	"ditto/booking/models"
	"ditto/booking/utils"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//CreateTenant -
func (d *Database) CreateTenant(db *gorm.DB, logon *cx.Payload, vmap map[string]interface{}) (*models.Tenant, error) {
	db = d.ValidDB(db)

	data := models.Tenant{}
	data.Name = vmap["name"].(string)
	data.UpdateUser = logon.ID
	result := db.Create(&data)

	if result.Error != nil {
		return nil, result.Error
	}
	if len(vmap) > 1 {
		//get values
		values := make([]string, 0)
		for k, v := range vmap {
			val := fmt.Sprintf("(%v,%q,\"%v\",%v)", data.ID, k, v, logon.ID)

			values = append(values, val)
		}

		//insert update
		sql := "insert into tenants_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
		err := db.Exec(sql).Error
		if err != nil {
			return nil, err
		}
	}

	return &data, nil
}

//GetTenant -
func (d *Database) GetTenant(db *gorm.DB, logon *cx.Payload, id int64, name string) (*models.TenantWithDetail, error) {
	db = d.ValidDB(db)

	//select tenant and details
	sql := "select t.*,d.detail as detail from tenants t left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from tenants_detail d group by id) d on (d.id=t.id)"
	sql += " where (t.id=? or ?) and (t.name like ? or ?)"

	ei := (id == 0)
	en := (name == "")
	data := models.TenantWithDetail{}
	result := db.Raw(sql, id, ei, fmt.Sprintf("%%%v%%", name), en).Scan(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == -1 {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}
	data.Detail = "{" + data.Detail + "}"

	err := utils.JSON.NewDecoder(strings.NewReader(data.Detail)).Decode(&data.Extra)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

//UpdateTenant -
func (d *Database) UpdateTenant(db *gorm.DB, logon *cx.Payload, tid int64, data map[string]interface{}) error {
	db = d.ValidDB(db)

	//update tenants
	name := data["name"].(string)
	if name != "" {
		sql := "update tenants set name=?,update_user=? where id=?"
		err := db.Exec(sql, name, logon.ID, tid).Error
		if err != nil {
			return err
		}
	}

	if len(data) > 1 {
		//get values
		values := make([]string, 0)
		for k, v := range data {
			val := fmt.Sprintf("(%v,%q,\"%v\",%v)", tid, k, v, logon.ID)

			values = append(values, val)
		}

		//insert update
		sql := "insert into tenants_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
		err := db.Exec(sql).Error
		if err != nil {
			return err
		}
	}

	return nil
}

//DeleteTenant -
func (d *Database) DeleteTenant(db *gorm.DB, logon *cx.Payload, tid int64) error {
	db = d.ValidDB(db)

	//delete tenants_detail
	sql := "delete from tenants_detail where id=?"
	err := db.Exec(sql, tid).Error
	if err != nil {
		return err
	}

	//delete tenants
	sql = "delete from tenants where id=?"
	err = db.Exec(sql, tid).Error
	if err != nil {
		return err
	}

	return nil
}

//AddUserToTenant -
func (d *Database) AddUserToTenant(db *gorm.DB, logon *cx.Payload, users []*models.TenantUsers) error {
	db = d.ValidDB(db)

	//clear
	sql := "update tenants_users set `right`=0 where user_id=? and `right`=1"
	err := db.Exec(sql, logon.ID).Error
	if err != nil {
		return err
	}

	//get values
	values := make([]string, 0)
	for _, u := range users {
		val := fmt.Sprintf("(%v,%v,%v,%v)", u.TenantID, u.UserID, u.Right, u.UpdateUser)

		values = append(values, val)
	}

	//insert update
	sql = "insert into tenants_users(`tenant_id`,`user_id`,`right`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update `right`=0,update_user=values(update_user)"
	err = db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

//GetTenantUsesr -
func (d *Database) GetTenantUsesr(db *gorm.DB, logon *cx.Payload, tid int64) ([]*models.TenantUsers, error) {
	db = d.ValidDB(db)

	//delete tenants_detail
	result := make([]*models.TenantUsers, 0)
	rows, err := db.Table("tenants_users").Where("tenant_id=?", tid).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.TenantUsers
		if err := db.ScanRows(rows, &record); err != nil {
			break
		}

		result = append(result, &record)
	}

	return result, nil
}

//GetUserTenants -
func (d *Database) GetUserTenants(db *gorm.DB, logon *cx.Payload, uid int64) ([]*models.TenantUsers, error) {
	db = d.ValidDB(db)
	//delete tenants_detail
	result := make([]*models.TenantUsers, 0)
	rows, err := db.Table("tenants_users").Where("user_id=?", uid).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.TenantUsers
		if err := db.ScanRows(rows, &record); err != nil {
			break
		}

		result = append(result, &record)
	}

	return result, nil
}

//ChangeUserTenant -
func (d *Database) ChangeUserTenant(db *gorm.DB, logon *cx.Payload, tid int64) error {
	db = d.ValidDB(db)

	//clear
	sql := "update tenants_users set `right`=0 where user_id=? and `right`=1"
	err := db.Exec(sql, logon.ID).Error
	if err != nil {
		return err
	}

	//
	sql = "update tenants_users set `right`=1 where user_id=? and `tenant_id`=?"
	err = db.Exec(sql, logon.ID, tid).Error
	if err != nil {
		return err
	}

	return nil
}
