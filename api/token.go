package api

import (
	"ditto/booking/config"
	"ditto/booking/cx"
	"ditto/booking/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//Token -
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Expires      int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) generateToken(d *cx.Payload) (*Token, error) {
	conf := config.Load()

	tm := time.Duration(conf.Expires)
	rftm := time.Duration(tm + 6)
	//create token
	token := jwt.New(jwt.SigningMethodRS512)

	b, err := utils.JSON.Marshal(d)
	if err != nil {
		return nil, err
	}
	//secret := security.EncryptSlice(b)
	secret := string(b)

	//set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = secret
	claims["exp"] = utils.NowJST().Time().Add(time.Hour * tm).Unix()

	//generate encoded token
	t, err := token.SignedString(s._rsa.GetPrivateKey())
	if err != nil {
		return nil, err
	}

	// secret = security.EncryptSlice(b)
	secret = string(b)
	expires := int64(time.Hour * rftm / time.Millisecond)
	//refresh token
	refreshtoken := jwt.New(jwt.SigningMethodRS512)
	rclaims := refreshtoken.Claims.(jwt.MapClaims)
	rclaims["sub"] = secret
	rclaims["exp"] = utils.NowJST().Time().Add(time.Hour * rftm).Unix()

	//generate encoded token
	rt, err := refreshtoken.SignedString(s._rsa.GetPrivateKey())
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  t,
		RefreshToken: rt,
		TokenType:    "Bearer",
		Expires:      expires,
	}, nil

}
