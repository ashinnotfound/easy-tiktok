package logic

import (
	"context"
	pb "easy-tiktok/apps/social/proto"
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

	return nil, nil
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
