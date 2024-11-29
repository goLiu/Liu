package models

import (
	"github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

type LoginInfo struct {
	IDNumber string `json:"id_number"` //身份证号码
	Password string `json:"password"`  //密码
}
