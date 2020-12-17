package db

import (
	"ditto/booking/models"
	"ditto/booking/utils"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

//UpdateUser -
func (d *Database) UpdateUser(db *gorm.DB, uid int64, email string, data map[string]interface{}) error {
	db = d.ValidDB(db)

	//get id
	sql := "select id from accounts where email = ?"

	var id int64
	result := db.Raw(sql, email).Scan(&id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == -1 {
		return gorm.ErrRecordNotFound
	}

	//get values
	values := make([]string, 0)
	for k, v := range data {
		val := fmt.Sprintf("(%v,%q,\"%v\",%v)", id, k, v, uid)

		values = append(values, val)
	}

	//insert update
	sql = "insert into users_detail(id,`option_key`,`option_val`,`update_user`) values " + strings.Join(values, ",") + " on duplicate key update option_val=values(option_val),update_user=values(update_user)"
	err := db.Exec(sql).Error
	if err != nil {
		return err
	}

	return nil
}

//UpdatePassword -
func (d *Database) UpdatePassword(db *gorm.DB, email, password string, upd time.Time) error {
	db = d.ValidDB(db)

	sql := "update accounts set password_hash=? where email=? and update_date=?"

	err := db.Exec(sql, password, email, upd).Error
	if err != nil {
		return err
	}
	return nil
}

//GetUser -
func (d *Database) GetUser(db *gorm.DB, email string) (*models.User, error) {
	db = d.ValidDB(db)

	//select user detail as a json format ("key":"value","key":"value")
	sql := "select d.id, d.detail from (select id, GROUP_CONCAT(CONCAT_WS(':', CONCAT('\"',`option_key`,'\"'), CONCAT('\"',`option_val`,'\"'))) as detail from users_detail d group by id) d left join accounts a on (a.id = d.id) where a.email = ?"

	var u struct {
		ID     int64
		Detail string
	}
	result := db.Raw(sql, email).Scan(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == -1 {
		return nil, gorm.ErrRecordNotFound
	}

	u.Detail = "{" + u.Detail + "}"
	data := models.User{
		ID:    u.ID,
		Email: email,
	}
	err := utils.JSON.NewDecoder(strings.NewReader(u.Detail)).Decode(&data)
	if err != nil {
		return nil, err
	}
	data.Email = email

	return &data, nil
}

//DeleteUser -
func (d *Database) DeleteUser(db *gorm.DB, email string) error {
	db = d.ValidDB(db)

	sql := "delete u from users_detail u left join accounts a on (a.id=u.id) where a.email=?"
	err := db.Exec(sql, email).Error
	if err != nil {
		return err
	}

	return nil
}
