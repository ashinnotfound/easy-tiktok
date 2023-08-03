package middleware

import (
	"easy-tiktok/util"
	"github.com/gin-gonic/gin"
)

func Jwt() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.URL.Path != "/douyin/user/login/" && context.Request.URL.Path != "/douyin/user/register/" && context.Request.URL.Path != "/douyin/feed" && context.Request.URL.Path != "/douyin/comment/list/" && context.Request.URL.Path != "/douyin/favorite/list/" && context.Request.URL.Path != "/douyin/publish/action/" {
			println(context.Request.URL.Path)
			value := context.Query("token")
			if value == "" {
				context.JSON(401, gin.H{
					"msg": "token is empty",
				})
				context.Abort()
				return
			}
			claims, err := util.ParseToken(value)
			if err != nil || claims.Secret != util.Key {
				context.JSON(401, gin.H{
					"msg": "token is invalid",
				})
				context.Abort()
				return
			}
		}
		context.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
