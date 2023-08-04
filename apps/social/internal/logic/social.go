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

	// 查询中间表
	var po model.UserFollow
	// 初始化响应
	response := &pb.DouyinRelationActionResponse{}
	response.StatusMsg = new(string)
	response.StatusCode = constant.RPC_STATUS.StatusOK()
	*response.StatusMsg = "RelationAction操作成功"

	var err error = nil
	// 开启事务
	tx := global.DB.Begin()
	// 逻辑判断
	result := tx.Where(&model.UserFollow{UserId: userId, FollowId: request.GetToUserId()}).Limit(1).Find(&po)
	if result.Error != nil {
		// 查询数据库出现问题
		response.StatusCode = constant.RPC_STATUS.StatusFailed()
		*response.StatusMsg = "数据库层面出现问题,RelationAction接口调用失败"
		global.LOGGER.Errorf("SocialServer::Action error: %v", result.Error)
		err = result.Error
	} else if result.RowsAffected == 0 && request.GetActionType() == constant.RELATION_NOT_FOLLOW {
		response.StatusCode = constant.RPC_STATUS.StatusFailed()
		*response.StatusMsg = "无效的输入参数,请注意前端代码"
		global.LOGGER.Warnf("SocialServer::Action error: %v", response.StatusMsg)
		err = errors.New(*response.StatusMsg)
	}

	if err == nil {
		// 数据库中没有记录
		if result.RowsAffected == 0 {
			po = model.UserFollow{UserId: userId, FollowId: request.GetToUserId(), Status: request.GetActionType()}
			tx.Create(&po)
		} else if po.Status != request.GetActionType() {
			po.Status = request.GetActionType()
			tx.Model(&po).Update("status", po.Status)
		}
		num := 1
		if request.GetActionType() == constant.RELATION_NOT_FOLLOW {
			num = -1
		}
		//操作user表
		if tx.Model(&orm.UserMsg{}).Where("id = ?", po.UserId).Update("follow_count", gorm.Expr("follow_count + ?", num)).Error != nil ||
			tx.Model(&orm.UserMsg{}).Where("id = ?", po.FollowId).Update("follower_count", gorm.Expr("follower_count + ?", num)).Error != nil {
			err = errors.New("操作user表失败")
			global.LOGGER.Errorf("SocialServer::Action error: %v", err)
			tx.Rollback()
		} else {
			tx.Commit()
		}
	} else {
		tx.Rollback()
	}
	return response, err
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
		global.LOGGER.Warnf("SocialServer::GetFollowList error: %v", result.Error)
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
					global.LOGGER.Warnf("SocialServer::GetFollowList error: %v", err)
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
		global.LOGGER.Warnf("SocialServer::GetFollowerList error: %v", result.Error)
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
					global.LOGGER.Warnf("SocialServer::GetFollowerList error: %v", err)
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
	// 初始化响应
	response := &pb.DouyinRelationFriendListResponse{}
	response.StatusMsg = new(string)
	*response.StatusMsg = "GetFriendList操作成功"
	response.StatusCode = constant.RPC_STATUS.StatusOK()

	var err error = nil

	// 通过用户id查询中间表
	var userFollowList []model.UserFollow
	result := global.DB.Where("user_id = ? AND status = ?", request.GetUserId(), constant.RELATION_FOLLOW).Find(&userFollowList)
	if result.Error != nil {
		err = result.Error
		*response.StatusMsg = "数据库层面出现问题,GetFriendList接口调用失败"
	} else if result.RowsAffected > 0 {
		wq := sync.WaitGroup{}

		friendIdChan := make(chan int64)

		wq.Add(1)
		// 需要判断是否相互关注
		go func() {
			defer wq.Done()
			for _, follow := range userFollowList {
				var friend model.UserFollow
				rowAffected := global.DB.Where("user_id = ? and follow_id = ? AND status = ?",
					follow.FollowId, request.GetUserId(), constant.RELATION_FOLLOW).Limit(1).Find(&friend).RowsAffected
				if rowAffected == 1 {
					friendIdChan <- friend.UserId
				}
			}
		}()

		wq.Add(1)
		go func(response *pb.DouyinRelationFriendListResponse) {
			// 用户的好友列表
			var friendList []*pb.FriendUser
			defer wq.Done()
			isFollow := true
			for {
				var friendMsg orm.UserMsg
				if userId, ok := <-friendIdChan; ok {
					// 获取好友信息
					global.DB.First(&friendMsg, userId)
					// 获取最新的聊天记录
					var messages [2]model.Message
					global.DB.Order("created_at asc").Where("from_user_id = ? AND to_user_id = ?", request.GetUserId(), userId).Limit(1).Find(messages[0])
					global.DB.Order("created_at asc").Where("from_user_id = ? AND to_user_id = ?", userId, request.GetUserId()).Limit(1).Find(messages[1])

					var message string
					var msgType int64
					// 比较两条消息的最新时间
					if messages[0].CreatedAt.Before(messages[1].CreatedAt) {
						message = messages[1].Content
						msgType = constant.MESSAGE_RECEIVE
					} else {
						message = messages[0].Content
						msgType = constant.MESSAGE_SEND
					}

					friend := &pb.User{
						Id:              &userId,
						Name:            &friendMsg.Username,
						FollowCount:     &friendMsg.FollowCount,
						FollowerCount:   &friendMsg.FollowerCount,
						IsFollow:        &isFollow,
						Avatar:          &friendMsg.Avatar.String,
						BackgroundImage: &friendMsg.BackgroundImage.String,
						Signature:       &friendMsg.Signature.String,
						TotalFavorited:  &friendMsg.TotalFavorited.Int64,
						WorkCount:       &friendMsg.WorkCount,
						FavoriteCount:   &friendMsg.FavoriteCount,
					}

					friendList = append(friendList, &pb.FriendUser{
						User:    friend,
						Message: &message,
						MsgType: &msgType,
					})
				} else {
					response.UserList = friendList
					break
				}

			}
		}(response)
		wq.Wait()
	}

	return response, err
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
