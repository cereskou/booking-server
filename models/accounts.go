package models

import (
	"time"
)

//Account -
type Account struct {
	ID                int64     `gorm:"column:id;primary_key"`                               //ID
	Email             string    `gorm:"column:email"`                                        //Email
	EmailConfirmed    int       `gorm:"column:email_confirmed"`                              //Email Confirmed
	PasswordHash      string    `gorm:"column:password_hash" json:"password_hash,omitempty"` // password hash
	LockoutEnd        time.Time `gorm:"column:lockout_end"`                                  //
	LockoutEnabled    int       `gorm:"column:lockout_enabled"`                              //
	AccessFailedCount int64     `gorm:"column:access_failed_count"`                          //
	LoginTime         time.Time `gorm:"column:login_time"`                                   //前回ログイン日時
	UpdateUser        int64     `gorm:"column:update_user"`                                  //更新者
	UpdateDate        time.Time `gorm:"column:update_date"`                                  //更新日時
}

//AccountWithRole -
type AccountWithRole struct {
	Account
	Name string `gorm:"column:name"` //名前
	Role string `gorm:"column:role"` //ユーザーロール
}

// TableName sets the insert table name for this struct type
func (a *Account) TableName() string {
	return "accounts"
}
