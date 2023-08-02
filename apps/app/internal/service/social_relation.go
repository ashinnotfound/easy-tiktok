package service

import (
	"context"
	"easy-tiktok/apps/app/internal/rpc"
	pb "easy-tiktok/apps/social/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RelationService //
// 关系方面的服务
// Author lql
type RelationService struct {
}

// Action //
// 关注or取消关注
func (service *RelationService) Action(c *gin.Context) {
	token := c.Query("token")
	toUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionTypeInt64, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)
	actionType := int32(actionTypeInt64)
	socialClient := rpc.GetSocialRpc()
	response, err := socialClient.RelationAction(context.Background(),
		&pb.DouyinRelationActionRequest{
			Token:      &token,
			ToUserId:   &toUserId,
			ActionType: &actionType,
		})
	if err != nil {
		//global.LOGGER.Errorf("RelationService::Action方法出错,reason: %v\n", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// FollowList //
// 获取关注列表
func (service *RelationService) FollowList(c *gin.Context) {

}

// FollowerList //
// 获取粉丝列表
func (service *RelationService) FollowerList(c *gin.Context) {

}

// FriendList //
// 获取好友列表
func (service *RelationService) FriendList(c *gin.Context) {

}
