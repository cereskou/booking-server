package models

import "time"

//Reservation -
type Reservation struct {
	ID         int64     `gorm:"column:id;primary_key" json:"id"`       //ID
	ScheduleID int64     `gorm:"column:schedule_id" json:"schedule_id"` //スケジュールID
	UserID     int64     `gorm:"column:user_id" json:"user_id"`         //ユーザーID
	Status     int       `gorm:"column:status" json:"status"`           //ステータス
	UpdateUser int64     `gorm:"column:update_user" json:"update_user"` //更新者
	UpdateDate time.Time `gorm:"column:update_date" json:"update_date"` //更新日時
}

//ReservationWithDetail -
type ReservationWithDetail struct {
	Reservation
	Detail string      `json:"-" gorm:"column:detail"` //詳細情報
	Extra  interface{} `json:"extra" gorm:"-"`
}

// TableName sets the insert table name for this struct type
func (r *Reservation) TableName() string {
	return "reservation"
}
