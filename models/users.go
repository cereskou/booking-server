package models

import (
	"time"
)

//User -
type User struct {
	ID         int64     `json:"id" gorm:"column:id;primary_key"`       //ID
	Email      string    `json:"email" gorm:"column:email"`             //Email
	LoginTime  time.Time `json:"login_time" gorm:"column:login_time"`   //前回ログイン日時
	UpdateUser int64     `json:"update_user" gorm:"column:update_user"` //更新者
	UpdateDate time.Time `json:"update_date" gorm:"column:update_date"` //更新日時
}

//UserWithDetail -
type UserWithDetail struct {
	User
	Detail string      `json:"-" gorm:"column:detail"` //詳細情報
	Extra  interface{} `json:"extra" gorm:"-"`
}

// //User -
// type User struct {
// 	ID         int64  `json:"id"`                        //ID
// 	Email      string `json:"email"`                     //Email
// 	Name       string `json:"name"`                      //氏名
// 	Contacts   string `json:"contact"`                   //連絡先
// 	Phone      string `json:"phone"`                     //電話番号
// 	Age        string `json:"age" type:"integer"`        //年齢
// 	Gender     string `json:"gender" type:"integer"`     //性別コード　0:男性 1:女性
// 	Occupation string `json:"occupation" type:"integer"` //職業コード
// }

//UserDetail -
type UserDetail struct {
	ID         int64     `json:"id" gorm:"column:id;primary_key"`                 //アカウントID
	OptionKey  string    `json:"option_key" gorm:"column:option_key;primary_key"` //属性
	OptionVal  string    `json:"option_val" gorm:"column:option_val"`             //値
	UpdateUser int64     `json:"update_user" gorm:"column:update_user"`           //更新者
	UpdateDate time.Time `json:"update_date" gorm:"column:update_date"`           //更新日時
}

// TableName sets the insert table name for this struct type
func (u *UserDetail) TableName() string {
	return "users_detail"
}
