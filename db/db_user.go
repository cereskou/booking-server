package db

import (
	"ditto/booking/cx"
	"ditto/booking/models"
	"ditto/booking/utils"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

//UpdateUser -
func (d *Database) UpdateUser(db *gorm.DB, logon *cx.Payload, id int64, data map[string]interface{}) error {
	db = d.ValidDB(db)

	// //get id
	// sql := "select id from accounts where email = ?"

	// var id int64
	// result := db.Raw(sql, email).Scan(&id)
	// if result.Error != nil {
	// 	return result.Error
	// }
	// if result.RowsAffected == -1 {
	// 	return gorm.ErrRecordNotFound
	// }

	//get values
	values := make([]string, 0)
	for k, v := range data {
		val := fmt.Sprintf("(%v,%q,\"%v\",%v)", id, k, v, logon.ID)

		values = append(values, val)
	}

	//insert update
	sql := "insert into users_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

//UpdatePassword -
func (d *Database) UpdatePassword(db *gorm.DB, logon *cx.Payload, password string, upd time.Time) error {
	db = d.ValidDB(db)

	sql := "update accounts set password_hash=?,update_user=? where email=? and update_date=?"

	err := db.Exec(sql, password, logon.ID, logon.Email, upd).Error
	if err != nil {
		return err
	}

	return nil
}

//GetUser -
func (d *Database) GetUser(db *gorm.DB, id int64) (*models.UserWithDetail, error) {
	db = d.ValidDB(db)

	//select user detail as a json format ("key":"value","key":"value")
	sql := "select a.id as id,a.email as email,a.login_time as login_time,a.update_user as update_user,a.update_date as update_date,d.detail as detail from accounts a"
	sql += " left join (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from users_detail d group by id) d on (a.id = d.id)"
	sql += " where a.id=?"

	var u models.UserWithDetail
	result := db.Raw(sql, id).Scan(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected <= 0 {
		return nil, gorm.ErrRecordNotFound
	}

	u.Detail = "{" + u.Detail + "}"
	err := utils.JSON.NewDecoder(strings.NewReader(u.Detail)).Decode(&u.Extra)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

//DeleteUser -
func (d *Database) DeleteUser(db *gorm.DB, id int64) error {
	db = d.ValidDB(db)

	sql := "delete u from users_detail u left join accounts a on (a.id=u.id) where a.id=?"
	err := db.Exec(sql, id).Error
	if err != nil {
		return err
	}

	return nil
}

//AddUserRole -
func (d *Database) AddUserRole(db *gorm.DB, logon *cx.Payload, uid int64, rids []int64) error {
	db = d.ValidDB(db)

	values := make([]string, 0)
	for _, id := range rids {
		val := fmt.Sprintf("(%v,%v,%v)", uid, id, logon.ID)
		values = append(values, val)
	}

	//insert update
	sql := "insert into accounts_roles(`account_id`,`role_id`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update update_user=values(update_user)"
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

//DeleteUserRole -
func (d *Database) DeleteUserRole(db *gorm.DB, logon *cx.Payload, uid int64, rids []int64) error {
	db = d.ValidDB(db)

	ids := make([]string, 0)
	for _, id := range rids {
		ids = append(ids, fmt.Sprintf("%v", id))
	}

	sql := "delete from accounts_roles where account_id=?"
	if len(ids) > 0 {
		sql += " and role_id in (?)"
	}
	err := db.Exec(sql, uid, strings.Join(ids, ",")).Error
	if err != nil {
		return err
	}

	return nil
}
