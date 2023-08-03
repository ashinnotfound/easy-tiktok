package logic

import (
	"context"
	"easy-tiktok/apps/constant"
	"easy-tiktok/apps/global"
	"easy-tiktok/apps/social/model"
	pb "easy-tiktok/apps/social/proto"
	orm "easy-tiktok/db/mysql"
	jwt "easy-tiktok/util"
	"errors"
	"gorm.io/gorm"
	"sync"
)

// SocialServerImpl //
// SocialServer接口的实现类
// Author lql
type SocialServerImpl struct {
	pb.SocialServer
}

// RelationAction //
// 取消关注和关注用户
func (impl *SocialServerImpl) RelationAction(ctx context.Context, request *pb.DouyinRelationActionRequest) (*pb.DouyinRelationActionResponse, error) {
	// 通过token获取用户id
	userId := jwt.GetUserId(request.GetToken())
	// 获取中间表
	userFollowTable := global.DB.Table(model.USER_FOLLOW_TABLE)
	// 查询中间表
	var po model.UserFollow
	// 初始化响应
	response := &pb.DouyinRelationActionResponse{}
	response.StatusMsg = new(string)
	response.StatusCode = constant.RPC_STATUS.StatusOK()

	// 逻辑判断
	if result := userFollowTable.Where(&model.UserFollow{UserId: userId, FollowId: request.GetToUserId()}).Limit(1).Find(&po); result.Error != nil {
		// 查询数据库出现问题
		response.StatusCode = constant.RPC_STATUS.StatusFailed()
		*response.StatusMsg = "数据库层面出现问题,RelationAction接口调用失败"
		global.LOGGER.Warnf("RelationServer::Action error: %v\n", result.Error)
		return response, result.Error
	} else if result.RowsAffected == 0 {
		// 数据库中没有记录
		if request.GetActionType() == constant.RELATION_FOLLOW {
			po = model.UserFollow{UserId: userId, FollowId: request.GetToUserId(), Status: request.GetActionType()}
			userFollowTable.Create(&po)
		} else if request.GetActionType() == constant.RELATION_NOT_FOLLOW {
			response.StatusCode = constant.RPC_STATUS.StatusFailed()
			*response.StatusMsg = "无效的输入参数,请注意前端代码"
			global.LOGGER.Warnf("RelationServer::Action error: %v\n", response.StatusMsg)
			return response, errors.New(*response.StatusMsg)
		}
	} else {
		if po.Status != request.GetActionType() {
			po.Status = request.GetActionType()
			userFollowTable.Model(&po).Update("status", po.Status)
		}
	}
	*response.StatusMsg = "RelationAction操作成功"
	return response, nil
}

// GetFollowList //
// 获取关注用户列表
func (impl *SocialServerImpl) GetFollowList(ctx context.Context, request *pb.DouyinRelationFollowListRequest) (*pb.DouyinRelationFollowListResponse, error) {
	// 初始化响应
	response := &pb.DouyinRelationFollowListResponse{}
	response.StatusMsg = new(string)
	response.StatusCode = constant.RPC_STATUS.StatusOK()
	wq := sync.WaitGroup{}
	// 初始化err
	var err error = nil
	// 声明userMsg
	userMsgChan := make(chan orm.UserMsg)
	var userFollowList []model.UserFollow
	if result := global.DB.Table(model.USER_FOLLOW_TABLE).Where("user_id = ? AND status = ?", request.GetUserId(), constant.RELATION_FOLLOW).Find(&userFollowList); result.Error != nil {
		response.StatusCode = constant.RPC_STATUS.StatusFailed()
		*response.StatusMsg = "数据库层面出现问题,GetFollowList接口调用失败"
		global.LOGGER.Warnf("RelationServer::GetFollowList error: %v", result.Error)
		return response, result.Error
	} else if result.RowsAffected > 0 {
		var followedUserList []*pb.User
		// 数据库读取
		wq.Add(1)
		go func(list []model.UserFollow) {
			defer wq.Done()
			// 关闭管道
			defer close(userMsgChan)
			for _, user := range userFollowList {
				var userMsg orm.UserMsg
				err = global.DB.First(&userMsg, user.FollowId).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					global.LOGGER.Warnf("RelationServer::GetFollowList error: %v", err)
				} else {
					userMsgChan <- userMsg
				}
			}
		}(userFollowList)
		wq.Add(1)
		// 添加到followedUserList
		go func(res *pb.DouyinRelationFollowListResponse) {
			defer wq.Done()
			status := true
			for {
				if userMsg, ok := <-userMsgChan; ok {
					followedUserList = append(followedUserList, &pb.User{
						Id:              &userMsg.ID,
						Name:            &userMsg.Username,
						FollowCount:     &userMsg.FollowCount,
						FollowerCount:   &userMsg.FollowerCount,
						IsFollow:        &status,
						Avatar:          &userMsg.Avatar.String,
						BackgroundImage: &userMsg.BackgroundImage.String,
						Signature:       &userMsg.Signature.String,
						TotalFavorited:  &userMsg.TotalFavorited.Int64,
						WorkCount:       &userMsg.WorkCount,
						FavoriteCount:   &userMsg.FavoriteCount,
					})
				} else {
					// 管道中的数据已经取完
					res.UserList = followedUserList
					break
				}
			}
		}(response)
	}

	wq.Wait()
	*response.StatusMsg = "GetFollowList操作成功"
	return response, err
}

// GetFollowerList //
// 获取登录用户粉丝列表
func (impl *SocialServerImpl) GetFollowerList(ctx context.Context, request *pb.DouyinRelationFollowerListRequest) (*pb.DouyinRelationFollowerListResponse, error) {
	// 初始化响应
	response := &pb.DouyinRelationFollowerListResponse{}
	response.StatusMsg = new(string)
	response.StatusCode = constant.RPC_STATUS.StatusOK()
	wq := sync.WaitGroup{}
	// 初始化err
	var err error = nil
	// 声明userMsg
	userMsgChan := make(chan orm.UserMsg)
	var userFollowerList []model.UserFollow
	if result := global.DB.Table(model.USER_FOLLOW_TABLE).Where("follow_id = ? AND status = ?", request.GetUserId(), constant.RELATION_FOLLOW).Find(&userFollowerList); result.Error != nil {
		response.StatusCode = constant.RPC_STATUS.StatusFailed()
		*response.StatusMsg = "数据库层面出现问题,GetFollowList接口调用失败"
		global.LOGGER.Warnf("RelationServer::GetFollowerList error: %v", result.Error)
		return response, result.Error
	} else if result.RowsAffected > 0 {
		var followedUserList []*pb.User
		// 数据库读取
		wq.Add(1)
		go func(list []model.UserFollow) {
			defer wq.Done()
			// 关闭管道
			defer close(userMsgChan)
			for _, user := range userFollowerList {
				var userMsg orm.UserMsg
				err = global.DB.First(&userMsg, user.UserId).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					global.LOGGER.Warnf("RelationServer::GetFollowerList error: %v", err)
				} else {
					userMsgChan <- userMsg
				}
			}
		}(userFollowerList)
		wq.Add(1)
		// 添加到followedUserList
		go func(res *pb.DouyinRelationFollowerListResponse) {
			defer wq.Done()
			status := true
			for {
				if userMsg, ok := <-userMsgChan; ok {
					followedUserList = append(followedUserList, &pb.User{
						Id:              &userMsg.ID,
						Name:            &userMsg.Username,
						FollowCount:     &userMsg.FollowCount,
						FollowerCount:   &userMsg.FollowerCount,
						IsFollow:        &status,
						Avatar:          &userMsg.Avatar.String,
						BackgroundImage: &userMsg.BackgroundImage.String,
						Signature:       &userMsg.Signature.String,
						TotalFavorited:  &userMsg.TotalFavorited.Int64,
						WorkCount:       &userMsg.WorkCount,
						FavoriteCount:   &userMsg.FavoriteCount,
					})
				} else {
					// 管道中的数据已经取完
					res.UserList = followedUserList
					break
				}
			}
		}(response)
	}

	wq.Wait()
	*response.StatusMsg = "GetFollowerList操作成功"
	return response, err
}

// GetFriendList //
// 获取登录用户好友列表
func (impl *SocialServerImpl) GetFriendList(ctx context.Context, request *pb.DouyinRelationFriendListRequest) (*pb.DouyinRelationFriendListResponse, error) {

	return nil, nil
}

// Chat //
// 获取对话消息
func (impl *SocialServerImpl) Chat(ctx context.Context, request *pb.DouyinMessageChatRequest) (*pb.DouyinMessageChatResponse, error) {
	// 显示全部（limit 10）

	return nil, nil
}

// MessageAction //
// 消息操作：发送消息
func (impl *SocialServerImpl) MessageAction(ctx context.Context, request *pb.DouyinMessageActionRequest) (*pb.DouyinMessageActionResponse, error) {
	var err error = nil
	// 初始化响应
	response := &pb.DouyinMessageActionResponse{}
	response.StatusCode = constant.RPC_STATUS.StatusOK()
	response.StatusMsg = new(string)

	// 判断是否为发送信息
	if request.GetActionType() == constant.MESSAGE_SEND {
		// 获取用户id
		userId := jwt.GetUserId(request.GetToken())

		// 将消息插入数据库
		if result := global.DB.Table(model.MESSAGE_TABLE).Create(&model.Message{
			FormUserID: userId,
			ToUserId:   request.GetToUserId(),
			Content:    request.GetContent(),
		}); result.Error != nil {
			response.StatusCode = constant.RPC_STATUS.StatusFailed()
			*response.StatusMsg = "数据库层面出现问题,MessageAction接口调用失败"
			global.LOGGER.Warn(response.StatusMsg)
			err = result.Error
		} else {
			*response.StatusMsg = "MessageAction操作成功"
		}
	} else {
		// 判断操作类型是否为发送消息
		response.StatusCode = constant.RPC_STATUS.StatusFailed()
		*response.StatusMsg = "无效的操作参数,请检查前端代码"
		global.LOGGER.Warn(response.StatusMsg)
		err = errors.New(*response.StatusMsg)
	}

	return response, err
}
