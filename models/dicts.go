package models

import (
	"time"
)

//Dict -
type Dict struct {
	ID         int64     `gorm:"column:id;primary_key"` //ID
	DictID     int       `gorm:"column:dict_id"`        //辞書ID
	Code       int       `gorm:"column:code"`           //コード
	Kvalue     string    `gorm:"column:kvalue"`         //値
	Remark     string    `gorm:"column:remark"`         //備考
	Status     int       `gorm:"column:status"`         //状態
	UpdateUser int64     `gorm:"column:update_user"`    //更新者
	UpdateDate time.Time `gorm:"column:update_date"`    //更新日時
}

// TableName sets the insert table name for this struct type
func (d *Dict) TableName() string {
	return "dicts"
}

//DictsType -
type DictsType struct {
	ID         int64     `gorm:"column:id;primary_key"` //ID
	DictID     int       `gorm:"column:dict_id"`        //辞書ID
	DictName   string    `gorm:"column:dict_name"`      //辞書名
	Remark     string    `gorm:"column:remark"`         //備考
	Status     int       `gorm:"column:status"`         //状態
	UpdateUser int64     `gorm:"column:update_user"`    //更新者
	UpdateDate time.Time `gorm:"column:update_date"`    //更新日時
}

// TableName sets the insert table name for this struct type
func (d *DictsType) TableName() string {
	return "dicts_type"
}
