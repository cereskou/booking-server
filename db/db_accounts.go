package db

import (
	"ditto/booking/models"
	"ditto/booking/utils"

	"gorm.io/gorm"
)

//GetAccount -
func (d *Database) GetAccount(email string) (*models.AccountWithRole, error) {
	data := models.AccountWithRole{}

	//select user and users_roles
	sql := "select u.*, d.option_val as name, s.role from accounts u left join (select ur.account_id, GROUP_CONCAT(r.name) as role from accounts_roles ur left join roles r on (ur.role_id = r.id) group by ur.account_id) s on (s.account_id = u.id) left join users_detail d on (u.id = d.id and d.`option_key`='name') where u.email = ?"

	result := d.DB().Raw(sql, email).Scan(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == -1 {
		return nil, gorm.ErrRecordNotFound
	}

	return &data, nil
}

//CreateAccount -
func (d *Database) CreateAccount(updid int64, vmap map[string]interface{}) (*models.Account, error) {
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
	result := d.DB().Create(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	//create a confirm_record

	return &data, nil
}

//CreateConfirmCode -
func (d *Database) CreateConfirmCode(data *models.Account) (*models.AccountConfirm, error) {
	code := utils.GeerateIDBase36()
	rec := models.AccountConfirm{
		AccountID:   data.ID,
		Email:       data.Email,
		ConfirmCode: code,
	}

	result := d.DB().Create(&rec)

	if result.Error != nil {
		return nil, result.Error
	}

	return &rec, nil
}

//GetConfirm -
func (d *Database) GetConfirm(email string, code string) (*models.AccountConfirm, error) {
	sql := "select * from accounts_confirm where confirm_code=? and email=?"

	rec := models.AccountConfirm{}
	result := d.DB().Raw(sql, code, email).Scan(&rec)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == -1 {
		return nil, gorm.ErrRecordNotFound
	}

	return &rec, nil
}

//DelConfirm -
func (d *Database) DelConfirm(id int64) error {
	sql := "delete from accounts_confirm where id=?"

	err := d.DB().Exec(sql, id).Error
	if err != nil {
		return err
	}

	return nil
}

//ConfirmAccount -
func (d *Database) ConfirmAccount(uid int64) error {
	sql := "update accounts set email_confirmed=1 where id=?"

	err := d.DB().Exec(sql, uid).Error
	if err != nil {
		return err
	}

	return nil
}

//ConfirmAccountWithCode -
func (d *Database) ConfirmAccountWithCode(email string, code string, expires int64) error {
	sql := "update accounts a,accounts_confirm ac set a.email_confirmed=1,ac.used = 1 where a.id=ac.account_id and ac.used=0 and ac.confirm_code=? and ac.email=?"
	//有効期限
	if expires > 0 {
		sql += " and TIME_TO_SEC(timediff(now(),ac.update_date))<=?"
	}

	result := d.DB().Exec(sql, code, email, expires)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected <= 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
