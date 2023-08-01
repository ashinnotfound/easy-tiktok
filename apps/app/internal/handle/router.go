package handle

import (
	"easy-tiktok/apps/app/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Route(e *gin.Engine) {
	e.Use(middleware.CORSMiddleware())
	e.Use(middleware.Jwt())
	group := e.Group("/douyin")
	{
		routerGroup := group.Group("/user")
		{
			routerGroup.POST("/login/", loginHandler)
			routerGroup.POST("/register/", registerHandler)
			routerGroup.GET("/", userInfoHandler)
		}
		group.GET("/feed", videoFeedHandler)
	}
}
