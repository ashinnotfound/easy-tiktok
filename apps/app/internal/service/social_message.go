package service

import (
	"context"
	"easy-tiktok/apps/app/internal/rpc"
	"easy-tiktok/apps/global"
	pb "easy-tiktok/apps/social/proto"
	"easy-tiktok/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// MessageService //
// 消息方面的服务
// Author lql
type MessageService struct {
	preRecord map[int64]preMsgTimeRecord
}

type preMsgTimeRecord struct {
	preTimestamp int64
	lock         sync.Mutex
}

func InitService() *MessageService {
	var service = &MessageService{}
	service.preRecord = make(map[int64]preMsgTimeRecord)
	return service
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
	// 更新最新时间
	userId := util.GetUserId(token)
	if record, ok := service.preRecord[userId]; ok {
		record.lock.Lock()
		record.preTimestamp = time.Now().Unix()
		record.lock.Unlock()
	}
	c.JSON(http.StatusOK, response)
}

// Chat //
// 聊天记录
func (service *MessageService) Chat(c *gin.Context) {
	// 获取请求参数
	token := c.Query("token")
	toUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	preMsgTime, _ := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64)
	userId := util.GetUserId(token)

	// 判断是否为第一次请求
	if preMsgTime == 0 {
		// 设置时间为当前请求的时间
		//requestTime := c.Request.Header.Get("Date")
		//parsedRequestTime, _ := time.Parse(time.RFC1123, requestTime)
		service.preRecord[userId] = preMsgTimeRecord{
			preTimestamp: time.Now().Unix(),
			lock:         sync.Mutex{},
		}
	} else {
		if timeRecord, ok := service.preRecord[userId]; ok {
			// 加锁
			timeRecord.lock.Lock()
			// 获取map中的时间
			preMsgTime = timeRecord.preTimestamp
			// 设置时间为当前请求的时间
			//requestTime := c.Request.Header.Get("Date")
			//parsedRequestTime, _ := time.Parse(time.RFC1123, requestTime)
			timeRecord.preTimestamp = time.Now().Unix()
			// 解锁
			timeRecord.lock.Unlock()
		}
	}
	global.LOGGER.Infof("pre_time=%v", preMsgTime)

	// 获取social-RPC客户端
	socialClient := rpc.GetSocialRpc()
	// 发送请求
	response, err := socialClient.Chat(context.Background(), &pb.DouyinMessageChatRequest{
		Token:      &token,
		ToUserId:   &toUserId,
		PreMsgTime: &preMsgTime,
	})
	if err != nil {
		global.LOGGER.Errorf("MessageService::Chat方法出错,reason: %v", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	c.JSON(http.StatusOK, response)
}
