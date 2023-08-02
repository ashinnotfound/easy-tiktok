package logic

import (
	"context"
	"easy-tiktok/apps/constant"
	"easy-tiktok/apps/global"
	"easy-tiktok/apps/social/model"
	pb "easy-tiktok/apps/social/proto"
	jwt "easy-tiktok/util"
	"errors"
)

// RelationServiceImpl //
// RelationService接口的实现类
// Author lql
type RelationServiceImpl struct {
	pb.RelationServer
}

// Action //
// 取消关注和关注用户
func (impl *RelationServiceImpl) Action(ctx context.Context, request *pb.DouyinRelationActionRequest) (*pb.DouyinRelationActionResponse, error) {
	// 通过token获取用户id
	userId := jwt.GetUserId(request.GetToken())

	// 获取中间表
	userFollowTable := global.DB.Table(model.USER_FOLLOW_TABLE)
	// 查询中间表
	var po model.UserFollow
	// 初始化响应
	var response *pb.DouyinRelationActionResponse
	*response.StatusCode = constant.STATUS_OK

	// 逻辑判断
	if result := userFollowTable.Where(&model.UserFollow{UserId: userId, FollowId: request.GetToUserId()}).Limit(1).Find(&po); result.Error != nil {
		// 查询数据库出现问题
		*response.StatusCode = constant.STATUS_FAILED
		*response.StatusMsg = "数据库层面出现问题,Action接口调用失败"
		global.LOGGER.Warnf("RelationServer::Action error: %v\n", response.StatusMsg)
		return response, result.Error
	} else if result.RowsAffected == 0 {
		// 数据库中没有记录
		if request.GetActionType() == 1 {
			po = model.UserFollow{UserId: userId, FollowId: request.GetToUserId(), Status: request.GetActionType()}
			userFollowTable.Create(&po)
		} else if request.GetActionType() == 2 {
			*response.StatusCode = constant.STATUS_FAILED
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
	*response.StatusMsg = "Action操作成功"
	return response, nil
}

// GetFollowList //
// 获取关注用户列表
func (impl *RelationServiceImpl) GetFollowList(ctx context.Context, request *pb.DouyinRelationFollowListRequest) (*pb.DouyinRelationFollowListResponse, error) {

	return nil, nil
}

// GetFollowerList //
// 获取登录用户粉丝列表
func (impl *RelationServiceImpl) GetFollowerList(ctx context.Context, request *pb.DouyinRelationFollowerListRequest) (*pb.DouyinRelationFollowerListResponse, error) {

	return nil, nil
}

// GetFriendList //
// 获取登录用户好友列表
func (impl *RelationServiceImpl) GetFriendList(ctx context.Context, request *pb.DouyinRelationFriendListRequest) (*pb.DouyinRelationFriendListResponse, error) {

	return nil, nil
}
