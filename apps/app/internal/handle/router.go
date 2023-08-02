package handle

import (
	"easy-tiktok/apps/app/internal/middleware"
	"easy-tiktok/apps/app/internal/router"
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
		favoriteGroup := group.Group("/favorite")
		{
			favoriteGroup.POST("/action/", favoriteHandler)
			favoriteGroup.GET("/list/", getFavoriteListHandler)
		}
		commentGroup := group.Group("/comment")
		{
			commentGroup.POST("/action/", commentHandler)
			commentGroup.GET("/list/", getCommentListHandler)
		}
		publishGroup := group.Group("/publish")
		{
			publishGroup.POST("/action", videoActionHandler)
			publishGroup.GET("/list", videoListHandler)
		}
	}

	socialRouter := router.SocialRouter{RootGroup: group}
	socialRouter.InitializeRouter()
}
