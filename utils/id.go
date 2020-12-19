package utils

import (
	"github.com/rs/xid"
)

//GenerateSpanID -
func GenerateSpanID() string {
	return xid.New().String()
}
