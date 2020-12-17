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
	Email    string `json:"email"`
	Password string `json:"password"`
}

//RefreshToken -
type RefreshToken struct {
	GrantType string `json:"grant_type"`
	Token     string `json:"refresh_token"`
}

//Payload - access token payload
type Payload struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
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
	DictID int64  `jsont:"dictid"`
	Code   int64  `jsont:"code"`
	Value  string `jsont:"value"`
	Remark string `jsont:"remark"`
	Status int    `jsont:"status"`
}

//Dicts -
type Dicts struct {
	Dict []*Dict `jsont:"dicts"`
}
