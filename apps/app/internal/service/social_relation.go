package service

import (
	"context"
	"easy-tiktok/apps/app/internal/rpc"
	"easy-tiktok/apps/global"
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

type friendListResponse struct {
	StatusCode *int32     `json:"status_code"`
	StatusMsg  *string    `json:"status_msg"`
	UserList   []*pb.User `json:"user_list"`
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
		global.LOGGER.Errorf("RelationService::Action方法出错,reason: %v\n", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// FollowList //
// 获取关注列表
func (service *RelationService) FollowList(c *gin.Context) {
	// 获取参数
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	token := c.Query("token")

	socialClient := rpc.GetSocialRpc()
	response, err := socialClient.GetFollowList(context.Background(),
		&pb.DouyinRelationFollowListRequest{
			UserId: &userId,
			Token:  &token,
		})
	if err != nil {
		global.LOGGER.Errorf("RelationService::FollowList方法出错,reason: %v\n", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// FollowerList //
// 获取粉丝列表
func (service *RelationService) FollowerList(c *gin.Context) {
	// 获取参数
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	token := c.Query("token")

	socialClient := rpc.GetSocialRpc()
	response, err := socialClient.GetFollowerList(context.Background(),
		&pb.DouyinRelationFollowerListRequest{
			UserId: &userId,
			Token:  &token,
		})
	if err != nil {
		global.LOGGER.Errorf("RelationService::FollowList方法出错,reason: %v\n", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// FriendList //
// 获取好友列表
func (service *RelationService) FriendList(c *gin.Context) {
	// 获取参数
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	token := c.Query("token")

	socialClient := rpc.GetSocialRpc()
	response, err := socialClient.GetFriendList(context.Background(),
		&pb.DouyinRelationFriendListRequest{
			UserId: &userId,
			Token:  &token,
		})
	var userList []*pb.User
	for _, friend := range response.UserList {
		userList = append(userList, friend.User)
	}
	newResponse := &friendListResponse{
		StatusCode: response.StatusCode,
		StatusMsg:  response.StatusMsg,
		UserList:   userList,
	}
	if err != nil {
		global.LOGGER.Errorf("RelationService::FollowList方法出错,reason: %v\n", err)
		c.JSON(http.StatusBadRequest, newResponse)
		return
	}
	c.JSON(http.StatusOK, newResponse)
}
