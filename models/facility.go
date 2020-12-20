package models

import "time"

//Facility - 施設
type Facility struct {
	ID         int64     `gorm:"column:id;primary_key" json:"id"`       //
	TenantID   int64     `gorm:"column:tenant_id" json:"tenant_id"`     //
	Name       string    `gorm:"column:name" json:"name"`               //名称
	Status     int       `gorm:"column:status" json:"status"`           //ステータス
	UpdateUser int64     `gorm:"column:update_user" json:"update_user"` //更新者
	UpdateDate time.Time `gorm:"column:update_date" json:"update_date"` //更新日時
}

// TableName sets the insert table name for this struct type
func (f *Facility) TableName() string {
	return "facility"
}

//FacilityWithDetail -
type FacilityWithDetail struct {
	Facility
	Detail string      `json:"-" gorm:"column:detail"` //詳細情報
	Extra  interface{} `json:"extra" gorm:"-"`
}

//FacilityDetail - 施設詳細
type FacilityDetail struct {
	ID         int64     `gorm:"column:id;primary_key" json:"id"`                 //施設トID
	OptionKey  string    `gorm:"column:option_key;primary_key" json:"option_key"` //属性
	OptionVal  string    `gorm:"column:option_val" json:"option_val"`             //値
	UpdateUser int64     `gorm:"column:update_user" json:"update_user"`           //更新者
	UpdateDate time.Time `gorm:"column:update_date" json:"update_date"`           //更新日時
}

// TableName sets the insert table name for this struct type
func (f *FacilityDetail) TableName() string {
	return "facility_detail"
}
