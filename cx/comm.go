package cx

//Payload - access token payload
type Payload struct {
	ID     int64  `json:"id"`    //ユーザーID
	Email  string `json:"email"` //メールアドレス
	Role   string `json:"role"`  //ロール
	Tenant int64  `json:"-"`     //テナントID
}
