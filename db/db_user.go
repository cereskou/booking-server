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
func (d *Database) GetUser(db *gorm.DB, id int64) (*models.User, error) {
	db = d.ValidDB(db)

	//select user detail as a json format ("key":"value","key":"value")
	sql := "select d.id, a.email, d.detail from (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from users_detail d group by id) d left join accounts a on (a.id = d.id) where a.id = ?"

	var u struct {
		ID     int64
		Email  string
		Detail string
	}
	result := db.Raw(sql, id).Scan(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == -1 {
		return nil, gorm.ErrRecordNotFound
	}

	u.Detail = "{" + u.Detail + "}"
	data := models.User{
		ID:    u.ID,
		Email: u.Email,
	}
	err := utils.JSON.NewDecoder(strings.NewReader(u.Detail)).Decode(&data)
	if err != nil {
		return nil, err
	}
	data.Email = u.Email

	return &data, nil
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
