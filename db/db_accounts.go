package db

import (
	"ditto/booking/models"
	"ditto/booking/utils"

	"gorm.io/gorm"
)

//GetAccount -
func (d *Database) GetAccount(db *gorm.DB, email string) (*models.AccountWithRole, error) {
	db = d.ValidDB(db)

	data := models.AccountWithRole{}

	//select user and users_roles
	sql := "select u.*, d.option_val as name, s.role from accounts u left join (select ur.account_id, GROUP_CONCAT(r.name) as role from accounts_roles ur left join roles r on (ur.role_id = r.id) group by ur.account_id) s on (s.account_id = u.id) left join users_detail d on (u.id = d.id and d.`option_key`='name') where u.email = ?"

	result := db.Raw(sql, email).Scan(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == -1 {
		return nil, gorm.ErrRecordNotFound
	}

	return &data, nil
}

//CreateAccount -
func (d *Database) CreateAccount(db *gorm.DB, updid int64, vmap map[string]interface{}) (*models.Account, error) {
	db = d.ValidDB(db)

	// input check
	// email
	// password
	// role
	// name
	// age
	// phone
	// contact
	// gender
	// occupation
	data := models.Account{}
	data.Email = vmap["email"].(string)
	data.PasswordHash = utils.HashPassword(vmap["password"].(string))
	data.UpdateUser = updid
	result := db.Create(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	//create a confirm_record

	return &data, nil
}

//CreateConfirmCode -
func (d *Database) CreateConfirmCode(db *gorm.DB, data *models.Account) (*models.AccountConfirm, error) {
	db = d.ValidDB(db)

	code := utils.GeerateIDBase36()
	rec := models.AccountConfirm{
		AccountID:   data.ID,
		Email:       data.Email,
		ConfirmCode: code,
	}

	result := db.Create(&rec)

	if result.Error != nil {
		return nil, result.Error
	}

	return &rec, nil
}

//GetConfirm -
func (d *Database) GetConfirm(db *gorm.DB, email string, code string) (*models.AccountConfirm, error) {
	db = d.ValidDB(db)

	sql := "select * from accounts_confirm where confirm_code=? and email=?"

	rec := models.AccountConfirm{}
	result := db.Raw(sql, code, email).Scan(&rec)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == -1 {
		return nil, gorm.ErrRecordNotFound
	}

	return &rec, nil
}

//DelConfirm -
func (d *Database) DelConfirm(db *gorm.DB, id int64) error {
	db = d.ValidDB(db)

	sql := "delete from accounts_confirm where id=?"

	err := db.Exec(sql, id).Error
	if err != nil {
		return err
	}

	return nil
}

//ConfirmAccount -
func (d *Database) ConfirmAccount(db *gorm.DB, uid int64) error {
	db = d.ValidDB(db)

	sql := "update accounts set email_confirmed=1 where id=?"

	err := db.Exec(sql, uid).Error
	if err != nil {
		return err
	}

	return nil
}

//ConfirmAccountWithCode -
func (d *Database) ConfirmAccountWithCode(db *gorm.DB, email string, code string, expires int64) error {
	db = d.ValidDB(db)

	sql := "update accounts a,accounts_confirm ac set a.email_confirmed=1,ac.used = 1 where a.id=ac.account_id and ac.used=0 and ac.confirm_code=? and ac.email=?"
	//有効期限
	if expires > 0 {
		sql += " and TIME_TO_SEC(timediff(now(),ac.update_date))<=?"
	}

	result := db.Exec(sql, code, email, expires)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

//DeleteAccount -
func (d *Database) DeleteAccount(db *gorm.DB, updid int64, email string) error {
	db = d.ValidDB(db)

	sql := "delete from accounts where email=?"
	err := db.Exec(sql, email).Error
	if err != nil {
		return err
	}

	return nil
}
