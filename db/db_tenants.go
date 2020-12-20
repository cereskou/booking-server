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
	if result.RowsAffected <= 0 {
		return nil, gorm.ErrRecordNotFound
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

//RemoveUserFromTenant -
func (d *Database) RemoveUserFromTenant(db *gorm.DB, logon *cx.Payload, tid int64, uids []int64) error {
	db = d.ValidDB(db)

	ids := make([]string, 0)
	for _, id := range uids {
		ids = append(ids, fmt.Sprintf("%v", id))
	}
	//delete
	sql := "delete from tenants_users where `tenant_id`=? and user_id in(?)"
	err := db.Exec(sql, tid, strings.Join(ids, ",")).Error
	if err != nil {
		return err
	}

	return nil
}

//GetTenantUserWithDetail - 複数ユーザー取得
func (d *Database) GetTenantUserWithDetail(db *gorm.DB, logon *cx.Payload, tid int64) ([]*models.UserWithDetail, error) {
	db = d.ValidDB(db)

	//select from accounts and users_detail
	sql := "select t.id as id,t.email as email,t.login_time as login_time,t.update_user as update_user,t.update_date as update_date,d.detail as detail from accounts t"
	sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from users_detail d group by id) d on (d.id=t.id)"
	sql += " left join tenants_users u on (u.user_id=t.id)"
	sql += " where u.tenant_id=?"

	result := make([]*models.UserWithDetail, 0)
	rows, err := db.Raw(sql, tid).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.UserWithDetail
		if err := db.ScanRows(rows, &record); err != nil {
			break
		}

		record.Detail = "{" + record.Detail + "}"
		err := utils.JSON.NewDecoder(strings.NewReader(record.Detail)).Decode(&record.Extra)
		if err != nil {
			return nil, err
		}

		result = append(result, &record)
	}
	if len(result) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return result, nil
}

//GetTenantUser -
func (d *Database) GetTenantUser(db *gorm.DB, logon *cx.Payload, tid int64) ([]*models.TenantUsers, error) {
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
	if len(result) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return result, nil
}

//GetUserTenant -
func (d *Database) GetUserTenant(db *gorm.DB, uid int64) (*models.TenantUsers, error) {
	db = d.ValidDB(db)

	var record models.TenantUsers
	result := db.Table(record.TableName()).
		Where("`user_id`=? and `right`=1", uid).
		First(&record)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &record, nil
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
	if len(result) == 0 {
		return nil, gorm.ErrRecordNotFound
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

//TenantCreateUser -
func (d *Database) TenantCreateUser(db *gorm.DB, logon *cx.Payload, vmap map[string]interface{}) (*models.Account, error) {
	db = d.ValidDB(db)

	data := models.Account{}
	data.Email = vmap["email"].(string)
	data.PasswordHash = utils.HashPassword(vmap["password"].(string))
	data.UpdateUser = logon.ID
	result := db.Create(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	//add user to tenant
	sql := "insert into tenants_users(`tenant_id`,`user_id`,`right`,`update_user`) values (?,?,?,?)"
	err := db.Exec(sql, logon.Tenant, data.ID, 1, logon.ID)
	if err.Error != nil {
		return nil, err.Error
	}

	return &data, nil
}

//DivideUserToTenant -
func (d *Database) DivideUserToTenant(db *gorm.DB, logon *cx.Payload, tid int64, uids []int64) error {
	db = d.ValidDB(db)

	//
	values := make([]string, 0)
	for _, u := range uids {
		val := fmt.Sprintf("(%v,%v,%v)", tid, u, logon.ID)

		values = append(values, val)
	}

	//insert update
	sql := "insert into tenants_users(`tenant_id`,`user_id`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update update_user=values(update_user)"
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}
