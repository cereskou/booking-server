package api

import (
	"time"
)

//Response -
type Response struct {
	Code  int         `json:"code"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// HTTPError represents an error that occurred while handling a request.
// fork echo.HTTPError for swagger
type HTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
}

//KeyValue -
type KeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

//Empty -
type Empty struct{}

//Login -
type Login struct {
	Email    string `json:"account"`
	Password string `json:"password"`
}

//RefreshToken -
type RefreshToken struct {
	GrantType string `json:"grant_type"`
	Token     string `json:"refresh_token"`
}

//User - update user password
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Password - update user password
type Password struct {
	Email       string    `json:"email"`        //Email
	NewPassword string    `json:"new_password"` //New Password
	OldPassword string    `json:"old_password"` //Old Password
	UpdateDate  time.Time `json:"update_date"`  //更新日時
}

//Dict -
type Dict struct {
	DictID int64  `json:"dict_id"`
	Code   int64  `json:"code"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
	Status int    `json:"status"`
}

//Dicts -
type Dicts struct {
	List []*Dict `json:"dicts"`
}

//DivideUser - divide exist user into tenants
type DivideUser struct {
	Divides int   `json:"divides"` //0: remove 1: add
	UserID  int64 `json:"user_id"`
}

//DivideUsers -
type DivideUsers struct {
	List []*DivideUser `json:"users"`
}
