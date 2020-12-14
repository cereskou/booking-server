package models

import "time"

//UsersRoles -
type UsersRoles struct {
	UserID     int64     `gorm:"column:user_id;primary_key"` //ユーザーID
	RoleID     int64     `gorm:"column:role_id;primary_key"` //ロールID
	UpdateUser int64     `gorm:"column:update_user"`         //更新者
	UpdateDate time.Time `gorm:"column:update_date"`         //更新日時
}

// TableName sets the insert table name for this struct type
func (u *UsersRoles) TableName() string {
	return "users_roles"
}
