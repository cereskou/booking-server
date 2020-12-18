package models

import (
	"time"
)

//MailTemplate -
type MailTemplate struct {
	ID         int64     `json:"id" gorm:"column:id;primary_key"`       //ID
	TenantID   int64     `json:"tenant_id" gorm:"column:tenant_id"`     //テナントID
	MailID     string    `json:"mail_id" gorm:"column:mail_id"`         //メールID
	Subject    string    `json:"subject" gorm:"column:subject"`         //サブジェクト
	Body       string    `json:"body" gorm:"column:body"`               //本文
	Enabled    int       `json:"enabled" gorm:"column:enabled"`         //有効・無効
	UpdateUser int64     `json:"update_user" gorm:"column:update_user"` //更新者
	UpdateDate time.Time `json:"update_date" gorm:"column:update_date"` //更新日時
}

// TableName sets the insert table name for this struct type
func (m *MailTemplate) TableName() string {
	return "mails_template"
}
