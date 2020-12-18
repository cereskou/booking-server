package models

import "time"

//AccountsRoles -
type AccountsRoles struct {
	AccountID  int64     `json:"account_id" gorm:"column:account_id;primary_key"` //アカウントID
	RoleID     int64     `json:"role_id" gorm:"column:role_id;primary_key"`       //ロールID
	UpdateUser int64     `json:"update_user" gorm:"column:update_user"`           //更新者
	UpdateDate time.Time `json:"update_date" gorm:"column:update_date"`           //更新日時
}

// TableName sets the insert table name for this struct type
func (a *AccountsRoles) TableName() string {
	return "accounts_roles"
}
