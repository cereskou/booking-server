package models

import "encoding/json"

//User -
type User struct {
	ID         int64       `json:"id"`                                  //ID
	Email      string      `json:"email,omitempty"`                     //Email
	Name       string      `json:"name,omitempty"`                      //氏名
	Contacts   string      `json:"contact,omitempty"`                   //連絡先
	Phone      string      `json:"phone,omitempty"`                     //電話番号
	Age        json.Number `json:"age,omitempty" type:"integer"`        //年齢
	Gender     json.Number `json:"gender,omitempty" type:"integer"`     //性別コード　0:男性 1:女性
	Occupation json.Number `json:"occupation,omitempty" type:"integer"` //職業コード
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
