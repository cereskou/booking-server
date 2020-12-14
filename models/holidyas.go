package models

import (
	"time"
)

//Holiday -
type Holiday struct {
	ID         int64     `gorm:"column:id;primary_key"` //
	Ymd        time.Time `gorm:"column:ymd"`            //祝日・休日月日
	Name       string    `gorm:"column:name"`           //祝日・休日名称
	Class      int       `gorm:"column:class"`          //種類（1:振替休日）
	UpdateUser int64     `gorm:"column:update_user"`    //更新者
	UpdateDate time.Time `gorm:"column:update_date"`    //更新日時
}

// TableName sets the insert table name for this struct type
func (h *Holiday) TableName() string {
	return "holidays"
}
