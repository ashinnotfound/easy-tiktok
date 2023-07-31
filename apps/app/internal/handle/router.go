package handle

import (
	"easy-tiktok/apps/app/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Route(e *gin.Engine) {
	e.Use(middleware.CORSMiddleware())
	e.POST("/douyin/user/login/", loginHandler)
	e.POST("/douyin/user/register/", registerHandler)
	e.GET("/douyin/feed", videoFeedHandler)
	e.Use(middleware.Jwt())
	e.GET("/douyin/user/", userInfoHandler)
}
