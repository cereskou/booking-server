package cx

//Payload - access token payload
type Payload struct {
	ID     int64  `json:"id"`     //ユーザーID
	Email  string `json:"email"`  //メールアドレス
	Name   string `json:"name"`   //名前
	Role   string `json:"role"`   //ロール
	Tenant int64  `json:"tenant"` //所属コンテナ
}
