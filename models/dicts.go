package models

import (
	"time"
)

//Dict -
type Dict struct {
	ID         int64     `json:"id" gorm:"column:id;primary_key"`       //ID
	TenantID   int64     `json:"tenant_id" gorm:"column:tenant_id"`     //テナントID
	DictID     int       `json:"dict_id" gorm:"column:dict_id"`         //辞書ID
	Code       int       `json:"code" gorm:"column:code"`               //コード
	Kvalue     string    `json:"kvalue" gorm:"column:kvalue"`           //値
	Remark     string    `json:"remark" gorm:"column:remark"`           //備考
	Status     int       `json:"status" gorm:"column:status"`           //状態
	UpdateUser int64     `json:"update_user" gorm:"column:update_user"` //更新者
	UpdateDate time.Time `json:"update_date" gorm:"column:update_date"` //更新日時
}

// TableName sets the insert table name for this struct type
func (d *Dict) TableName() string {
	return "dicts"
}

//DictsType -
type DictsType struct {
	ID         int64     `json:"id" gorm:"column:id;primary_key"`       //ID
	DictID     int       `json:"dict_id" gorm:"column:dict_id"`         //辞書ID
	DictName   string    `json:"dict_name" gorm:"column:dict_name"`     //辞書名
	Remark     string    `json:"remark" gorm:"column:remark"`           //備考
	Status     int       `json:"status" gorm:"column:status"`           //状態
	UpdateUser int64     `json:"update_user" gorm:"column:update_user"` //更新者
	UpdateDate time.Time `json:"update_date" gorm:"column:update_date"` //更新日時
}

// TableName sets the insert table name for this struct type
func (d *DictsType) TableName() string {
	return "dicts_type"
}
