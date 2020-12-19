package models

import "time"

//Role -
type Role struct {
	ID         int64     `json:"id" gorm:"column:id;primary_key"`       //Id
	Name       string    `json:"name" gorm:"column:name"`               //ロール名
	UpdateUser int64     `json:"update_user" gorm:"column:update_user"` //更新者
	UpdateDate time.Time `json:"update_date" gorm:"column:update_date"` //更新日時
}

// TableName sets the insert table name for this struct type
func (r *Role) TableName() string {
	return "roles"
}
