package db

import (
	"ditto/booking/models"

	"gorm.io/gorm"
)

//GetUser -
func (d *Database) GetUser(email string) (*models.UserWithRole, error) {
	data := models.UserWithRole{}

	//select user and users_roles
	sql := "select u.*, r.role from users u left join (select ur.user_id, GROUP_CONCAT(r0.name) as role from users_roles ur left join roles r0 on (ur.role_id = r0.id) group by ur.user_id) r on (r.user_id = u.id) where u.email = ?"

	result := d.DB().Raw(sql, email).Scan(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == -1 {
		return nil, gorm.ErrRecordNotFound
	}

	return &data, nil
}
