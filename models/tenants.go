package models

import (
	"time"
)

//Tenant -
type Tenant struct {
	ID         int64     `json:"id" gorm:"column:id;primary_key"`       //テナントID
	Name       string    `json:"name" gorm:"column:name"`               //名前
	UpdateUser int64     `json:"update_user" gorm:"column:update_user"` //更新者
	UpdateDate time.Time `json:"update_date" gorm:"column:update_date"` //更新日時
}

// TableName sets the insert table name for this struct type
func (t *Tenant) TableName() string {
	return "tenants"
}

//TenantWithDetail -
type TenantWithDetail struct {
	Tenant
	Detail string      `json:"-" gorm:"column:detail"` //詳細情報
	Extra  interface{} `json:"extra" gorm:"-"`
}

//TenantDetail -
type TenantDetail struct {
	ID         int64     `json:"id" gorm:"column:id;primary_key"`                 //テナントID
	OptionKey  string    `json:"option_key" gorm:"column:option_key;primary_key"` //属性
	OptionVal  string    `json:"option_val" gorm:"column:option_val"`             //値
	UpdateUser int64     `json:"update_user" gorm:"column:update_user"`           //更新者
	UpdateDate time.Time `json:"update_date" gorm:"column:update_date"`           //更新日時
}

// TableName sets the insert table name for this struct type
func (t *TenantDetail) TableName() string {
	return "tenants_detail"
}

//TenantUsers -
type TenantUsers struct {
	TenantID   int64     `json:"tenant_id" gorm:"column:tenant_id;primary_key"` //テナントID
	UserID     int64     `json:"user_id" gorm:"column:user_id;primary_key"`     //ユーザーID
	Right      int       `json:"right" gorm:"column:right"`                     //権限
	UpdateUser int64     `json:"update_user" gorm:"column:update_user"`         //更新者
	UpdateDate time.Time `json:"update_date" gorm:"column:update_date"`         //更新日時
}

// TableName sets the insert table name for this struct type
func (t *TenantUsers) TableName() string {
	return "tenants_users"
}
