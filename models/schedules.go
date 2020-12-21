package models

import (
	"time"
)

//Schedule -
type Schedule struct {
	ID         int64     `gorm:"column:id;primary_key" json:"id"`       //ID
	MenuID     int64     `gorm:"column:menu_id" json:"menu_id"`         //メニューID
	YmdStart   time.Time `gorm:"column:ymd_start" json:"ymd_start"`     //開始日
	YmdEnd     time.Time `gorm:"column:ymd_end" json:"ymd_end"`         //終了日
	TimeStart  time.Time `gorm:"column:time_start" json:"time_start"`   //開始時間
	TimeEnd    time.Time `gorm:"column:time_end" json:"time_end"`       //終了時間
	Repetition int       `gorm:"column:repetition" json:"repetition"`   //繰り返し
	Capacity   int       `gorm:"column:capacity" json:"capacity"`       //定員数
	FacilityID int       `gorm:"column:facility_id" json:"facility_id"` //施設ID
	Status     int       `gorm:"column:status" json:"status"`           //ステータス
	UpdateUser int64     `gorm:"column:update_user" json:"update_user"` //更新者
	UpdateDate time.Time `gorm:"column:update_date" json:"update_date"` //更新日時
}

// TableName sets the insert table name for this struct type
func (s *Schedule) TableName() string {
	return "schedules"
}
