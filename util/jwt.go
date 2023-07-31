package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type DouYinClaims struct {
	UserId int64  `json:"userId"`
	Secret string `json:"secret"`
	jwt.StandardClaims
}

var Key = "JustAsNoStarsCanBeSeen"
var times int64 = 3600 * 24 * 7

func GetToken(userId int64) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, DouYinClaims{
		UserId: userId,
		Secret: Key,
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: time.Now().Unix() + times,
			Id:        "",
			IssuedAt:  time.Now().Unix(),
			Issuer:    "admin",
			NotBefore: 0,
			Subject:   "",
		},
	})
	signedString, err := claims.SignedString([]byte("golang"))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func ParseToken(token string) (*DouYinClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &DouYinClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("golang"), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*DouYinClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
func GetTokenTid(token string) int64 {
	claims, _ := ParseToken(token)
	return claims.UserId
}
