package cmd

import (
	"ditto/booking/config"
	"ditto/booking/db"
	"ditto/booking/models"
	"ditto/booking/utils"
	"strings"
)

//Migrate -
func Migrate(db *db.Database) error {
	conf := config.Load()

	mapRoles := make(map[string]int64)
	rows, _ := db.DB().Table("roles").Rows()
	defer rows.Close()
	for rows.Next() {
		role := models.Roles{}
		db.DB().ScanRows(rows, &role)

		mapRoles[role.Name] = role.ID
	}

	for _, u := range conf.Account {
		user := models.User{
			Email:        u.Email,
			PasswordHash: utils.HashPassword(u.Password),
		}

		result := db.DB().Table("users").Where("email=?", user.Email).Updates(map[string]interface{}{"password_hash": user.PasswordHash})
		if result.Error != nil || result.RowsAffected == 0 {
			result = db.DB().Table("users").Create(&user)
			if result.Error != nil {
				return result.Error
			}
		}
		//select user
		err := db.DB().Table("users").Where("email=?", user.Email).Scan(&user).Error
		if err != nil {
			return err
		}
		roles := strings.Split(u.Role, ",")
		for _, role := range roles {
			if id, ok := mapRoles[role]; ok {
				ur := models.UsersRoles{
					UserID: user.ID,
					RoleID: id,
				}
				err := db.DB().Table(ur.TableName()).Create(&ur).Error
				if err != nil {
					// return err
				}
			}
		}
	}

	return nil
}
