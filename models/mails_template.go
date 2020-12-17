package models

import (
	"time"
)

//MailTemplate -
type MailTemplate struct {
	ID         int64     `gorm:"column:id;primary_key"` //ID
	TenantID   int64     `gorm:"column:tenant_id"`      //テナントID
	MailID     string    `gorm:"column:mail_id"`        //メールID
	Subject    string    `gorm:"column:subject"`        //サブジェクト
	Body       string    `gorm:"column:body"`           //本文
	Enabled    int       `gorm:"column:enabled"`        //有効・無効
	UpdateUser int64     `gorm:"column:update_user"`    //更新者
	UpdateDate time.Time `gorm:"column:update_date"`    //更新日時
}

// TableName sets the insert table name for this struct type
func (m *MailTemplate) TableName() string {
	return "mails_template"
}
