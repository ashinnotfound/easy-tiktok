package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

type DouYinClaims struct {
	UserId int64  `json:"userId"`
	Secret string `json:"secret"`
	jwt.StandardClaims
}

var Key = "JustAsNoStarsCanBeSeen"

func Jwt() gin.HandlerFunc {
	return func(context *gin.Context) {
		value := context.Param("token")
		if value == "" {
			context.JSON(200, gin.H{
				"code": 500,
				"msg":  "token is empty",
			})
			context.Abort()
			return
		}
		claims, err := ParseToken(value)
		if err != nil || claims.Secret != Key {
			context.JSON(200, gin.H{
				"code": 500,
				"msg":  "token is invalid",
			})
			context.Abort()
			return
		}
		context.Next()
	}
}

func GetToken(userId int64, times int64) (string, error) {
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
