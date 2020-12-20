package db

import (
	"ditto/booking/cx"
	"ditto/booking/models"
	"ditto/booking/utils"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//CreateClass -
func (d *Database) CreateClass(db *gorm.DB, logon *cx.Payload, vmap map[string]interface{}) (*models.Class, error) {
	db = d.ValidDB(db)

	data := models.Class{}
	data.Name = vmap["name"].(string)
	data.TenantID = logon.Tenant
	data.OwnerID = logon.ID
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
		sql := "insert into classes_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
		err := db.Exec(sql).Error
		if err != nil {
			return nil, err
		}
	}

	return &data, nil
}

//GetClass - single
func (d *Database) GetClass(db *gorm.DB, logon *cx.Payload, id int64, name string) (*models.ClassWithDetail, error) {
	db = d.ValidDB(db)

	//get class and details
	sql := "select t.*,d.detail as detail from classes t left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from classes_detail d group by id) d on (d.id=t.id)"
	sql += " where (t.id=? or ?) and (t.name like ? or ?) and t.tenant_id=?"

	ei := (id == 0)
	en := (name == "")

	data := models.ClassWithDetail{}
	result := db.Raw(sql, id, ei, fmt.Sprintf("%%%v%%", name), en, logon.Tenant).Scan(&data)
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

//UpdateClass -
func (d *Database) UpdateClass(db *gorm.DB, logon *cx.Payload, cid int64, data map[string]interface{}) error {
	db = d.ValidDB(db)

	//update tenants
	name := data["name"].(string)
	if name != "" {
		sql := "update classes set name=?,update_user=? where id=?"
		err := db.Exec(sql, name, logon.ID, cid).Error
		if err != nil {
			return err
		}
	}

	if len(data) > 0 {
		//get values
		values := make([]string, 0)
		for k, v := range data {
			val := fmt.Sprintf("(%v,%q,\"%v\",%v)", cid, k, v, logon.ID)

			values = append(values, val)
		}

		//insert update
		sql := "insert into classes_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
		err := db.Exec(sql).Error
		if err != nil {
			return err
		}
	}

	return nil
}

//DeleteClass -
func (d *Database) DeleteClass(db *gorm.DB, logon *cx.Payload, cid int64) error {
	db = d.ValidDB(db)

	//delete classes_detail
	sql := "delete from classes_detail where id=?"
	err := db.Exec(sql, cid).Error
	if err != nil {
		return err
	}

	//delete classes
	sql = "delete from classes where id=?"
	err = db.Exec(sql, cid).Error
	if err != nil {
		return err
	}

	return nil
}

//AddUserToClass -
func (d *Database) AddUserToClass(db *gorm.DB, logon *cx.Payload, users []*models.ClassesUser) error {
	db = d.ValidDB(db)

	//get values
	values := make([]string, 0)
	for _, u := range users {
		val := fmt.Sprintf("(%v,%v,%v,%v)", u.ClassID, u.UserID, u.Status, u.UpdateUser)

		values = append(values, val)
	}

	//insert update
	sql := "insert into classes_users(`class_id`,`user_id`,`status`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update update_user=values(update_user)"
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

//RemoveUserFromClass -
func (d *Database) RemoveUserFromClass(db *gorm.DB, logon *cx.Payload, cid int64, uids []int64) error {
	db = d.ValidDB(db)

	ids := make([]string, 0)
	for _, id := range uids {
		ids = append(ids, fmt.Sprintf("%v", id))
	}
	//delete
	sql := "delete from classes_users where `class_id`=? and user_id in(?)"
	err := db.Exec(sql, cid, strings.Join(ids, ",")).Error
	if err != nil {
		return err
	}

	return nil
}

//GetClassUsersWithDetail - 複数ユーザー取得
func (d *Database) GetClassUsersWithDetail(db *gorm.DB, logon *cx.Payload, cid int64) ([]*models.UserWithDetail, error) {
	db = d.ValidDB(db)

	//select from accounts and users_detail
	sql := "select t.id as id,t.email as email,t.login_time as login_time,t.update_user as update_user,t.update_date as update_date,d.detail as detail from accounts t"
	sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from users_detail d group by id) d on (d.id=t.id)"
	sql += " left join classes_users u on (u.user_id=t.id)"
	sql += " where u.class_id=?"

	result := make([]*models.UserWithDetail, 0)
	rows, err := db.Raw(sql, cid).Rows()
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

//GetClassUsers -
func (d *Database) GetClassUsers(db *gorm.DB, logon *cx.Payload, cid int64) ([]*models.ClassesUser, error) {
	db = d.ValidDB(db)

	//delete tenants_detail
	result := make([]*models.ClassesUser, 0)
	rows, err := db.Table("classes_users").Where("class_id=?", cid).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.ClassesUser
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

//ChangeUserClass - ログイン中のユーザー
func (d *Database) ChangeUserClass(db *gorm.DB, logon *cx.Payload, tid int64) error {
	db = d.ValidDB(db)

	//clear
	sql := "update classes_users set `status`=0 where user_id=? and `status`=1"
	err := db.Exec(sql, logon.ID).Error
	if err != nil {
		return err
	}

	//
	sql = "update classes_users set `right`=1 where user_id=? and `tenant_id`=?"
	err = db.Exec(sql, logon.ID, tid).Error
	if err != nil {
		return err
	}

	return nil
}

//ClassCreateUser -
func (d *Database) ClassCreateUser(db *gorm.DB, logon *cx.Payload, cid int64, vmap map[string]interface{}) (*models.Account, error) {
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

	//add user to class
	sql = "insert into classes_users(`class_id`,`user_id`,`status`,`update_user`) values (?,?,?,?)"
	err = db.Exec(sql, cid, data.ID, 1, logon.ID)
	if err.Error != nil {
		return nil, err.Error
	}

	return &data, nil
}

//DivideUserToClass -
func (d *Database) DivideUserToClass(db *gorm.DB, logon *cx.Payload, cid int64, uids []int64) error {
	db = d.ValidDB(db)

	//
	values := make([]string, 0)
	for _, u := range uids {
		val := fmt.Sprintf("(%v,%v,%v)", cid, u, logon.ID)

		values = append(values, val)
	}

	//insert update
	sql := "insert into classes_users(`class_id`,`user_id`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update update_user=values(update_user)"
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

//GetUserClasses -
func (d *Database) GetUserClasses(db *gorm.DB, logon *cx.Payload, uid int64) ([]*models.ClassesUser, error) {
	db = d.ValidDB(db)

	result := make([]*models.ClassesUser, 0)
	rows, err := db.Table("classes_users").Where("user_id=?", uid).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record models.ClassesUser
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
