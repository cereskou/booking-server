package models

import (
	"time"
)

//Account -
type Account struct {
	ID                int64     `json:"id" gorm:"column:id;primary_key"`                       //ID
	Email             string    `json:"email" gorm:"column:email"`                             //Email
	EmailConfirmed    int       `json:"email_confirmed" gorm:"column:email_confirmed"`         //Email Confirmed
	PasswordHash      string    `json:"password_hash,omitempty" gorm:"column:password_hash"`   // password hash
	LockoutEnd        time.Time `json:"lockout_end" gorm:"column:lockout_end"`                 //
	LockoutEnabled    int       `json:"lockout_enabled" gorm:"column:lockout_enabled"`         //
	AccessFailedCount int64     `json:"access_failed_count" gorm:"column:access_failed_count"` //
	LoginTime         time.Time `json:"login_time" gorm:"column:login_time"`                   //前回ログイン日時
	UpdateUser        int64     `json:"update_user" gorm:"column:update_user"`                 //更新者
	UpdateDate        time.Time `json:"update_date" gorm:"column:update_date"`                 //更新日時
}

//AccountWithRole -
type AccountWithRole struct {
	Account
	Name     string `json:"name" gorm:"column:name"`        //名前
	Role     string `json:"role" gorm:"column:role"`        //ユーザーロール
	TenantID int64  `json:"tenant_id" gorm:"column:tenant"` //所属コンテナ
}

// TableName sets the insert table name for this struct type
func (a *Account) TableName() string {
	return "accounts"
}
