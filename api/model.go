package api

//Response -
type Response struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

// HTTPError represents an error that occurred while handling a request.
// fork echo.HTTPError for swagger
type HTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
}

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
