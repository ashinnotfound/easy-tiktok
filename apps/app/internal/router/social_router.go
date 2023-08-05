package router

import (
	"easy-tiktok/apps/app/internal/service"
	"easy-tiktok/apps/global"
	"github.com/gin-gonic/gin"
)

// SocialRouter //
// 社交路由
// Author lql
type SocialRouter struct {
	RootGroup *gin.RouterGroup
}

// InitializeRouter //
// 初始化路由
func (router *SocialRouter) InitializeRouter() {
	if router.RootGroup == nil {
		global.LOGGER.Error("social router not initialize")
	}
	var relationService service.RelationService
	messageService := service.InitService()
	relationGroup := router.RootGroup.Group("/relation")
	{
		relationGroup.POST("/action/", relationService.Action)
		relationGroup.GET("/follow/list/", relationService.FollowList)
		relationGroup.GET("/follower/list/", relationService.FollowerList)
		relationGroup.GET("/friend/list/", relationService.FriendList)
	}
	messageGroup := router.RootGroup.Group("/message")
	{
		messageGroup.POST("/action/", messageService.Action)
		messageGroup.GET("/chat/", messageService.Chat)
	}
}
