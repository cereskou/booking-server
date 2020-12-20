package models

import "time"

//Class -
type Class struct {
	ID         int64     `gorm:"column:id;primary_key" json:"id"`       //ID
	TenantID   int64     `gorm:"column:tenant_id" json:"tenant_id"`     //テナントID
	OwnerID    int64     `gorm:"column:owner_id" json:"owner_id"`       //担任ID
	Name       string    `gorm:"column:name" json:"name"`               //名前
	UpdateUser int64     `gorm:"column:update_user" json:"update_user"` //更新者
	UpdateDate time.Time `gorm:"column:update_date" json:"update_date"` //更新日時
}

//ClassWithDetail -
type ClassWithDetail struct {
	Class
	Detail string      `json:"-" gorm:"column:detail"` //詳細情報
	Extra  interface{} `json:"extra" gorm:"-"`
}

// TableName sets the insert table name for this struct type
func (c *Class) TableName() string {
	return "classes"
}

//ClassesDetail -
type ClassesDetail struct {
	ID         int64     `gorm:"column:id;primary_key" json:"id"`                 //クラスID
	OptionKey  string    `gorm:"column:option_key;primary_key" json:"option_key"` //属性
	OptionVal  string    `gorm:"column:option_val" json:"option_val"`             //値
	UpdateUser int64     `gorm:"column:update_user" json:"update_user"`           //更新者
	UpdateDate time.Time `gorm:"column:update_date" json:"update_date"`           //更新日時}
}

// TableName sets the insert table name for this struct type
func (c *ClassesDetail) TableName() string {
	return "classes_detail"
}

//ClassesUser -
type ClassesUser struct {
	ClassID    int64     `gorm:"column:class_id;primary_key" json:"class_id"` //クラストID
	UserID     int64     `gorm:"column:user_id;primary_key" json:"user_id"`   //ユーザーID
	Status     int       `gorm:"column:status" json:"status"`                 //ステータス
	UpdateUser int64     `gorm:"column:update_user" json:"update_user"`       //更新者
	UpdateDate time.Time `gorm:"column:update_date" json:"update_date"`       //更新日時
}

// TableName sets the insert table name for this struct type
func (c *ClassesUser) TableName() string {
	return "classes_users"
}
