package models

import "time"

//Roles -
type Roles struct {
	ID         int64     `gorm:"column:id;primary_key"` //Id
	Name       string    `gorm:"column:name"`           //ロール名
	UpdateUser int64     `gorm:"column:update_user"`    //更新者
	UpdateDate time.Time `gorm:"column:update_date"`    //更新日時
}

// TableName sets the insert table name for this struct type
func (r *Roles) TableName() string {
	return "roles"
}
