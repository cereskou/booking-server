package db

import (
	"ditto/booking/models"

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
