package models

import (
	"time"
)

//User -
type User struct {
	ID                int64     `gorm:"column:id;primary_key"`         //ID
	Name              string    `gorm:"column:name"`                   //名前
	Email             string    `gorm:"column:email"`                  //Email
	EmailConfirmed    bool      `gorm:"column:email_confirmed"`        //Email Confirmed
	PasswordHash      string    `gorm:"column:password_hash" json:"-"` //
	LockoutEnd        time.Time `gorm:"column:lockout_end"`            //
	LockoutEnabled    bool      `gorm:"column:lockout_enabled"`        //
	AccessFailedCount int       `gorm:"column:access_failed_count"`    //
	LoginTime         time.Time `gorm:"column:login_time"`             //前回ログイン日時
	UpdateUser        int64     `gorm:"column:update_user"`            //更新者
	UpdateDate        time.Time `gorm:"column:update_date"`            //更新日時
}

//UserWithRole -
type UserWithRole struct {
	User
	Role string `gorm:"column:role"` //ユーザーロール
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "users"
}
