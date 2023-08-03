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

// MessageService //
// 消息方面的服务
// Author lql
type MessageService struct {
}

// Action //
// 发送消息
func (service *MessageService) Action(c *gin.Context) {
	// 获取请求参数
	token := c.Query("token")
	toUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionTypeInt64, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)
	actionType := int32(actionTypeInt64)
	content := c.Query("content")

	// 获取social-RPC客户端
	socialClient := rpc.GetSocialRpc()

	// 发送请求
	response, err := socialClient.MessageAction(context.Background(), &pb.DouyinMessageActionRequest{
		Token:      &token,
		ToUserId:   &toUserId,
		Content:    &content,
		ActionType: &actionType,
	})
	if err != nil {
		global.LOGGER.Errorf("MessageService::Action方法出错,reason: %v", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

// Chat //
// 聊天记录
func (service *MessageService) Chat(c *gin.Context) {
	// 获取请求参数
	token := c.Query("token")
	toUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)

	// 获取social-RPC客户端
	socialClient := rpc.GetSocialRpc()

	// 发送请求
	response, err := socialClient.Chat(context.Background(), &pb.DouyinMessageChatRequest{
		Token:    &token,
		ToUserId: &toUserId,
	})
	if err != nil {
		global.LOGGER.Errorf("MessageService::Action方法出错,reason: %v", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, response)
}
