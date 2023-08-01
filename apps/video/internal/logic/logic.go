package logic

import (
	proto2 "easy-tiktok/apps/user/proto"
	"easy-tiktok/apps/video/proto"
	Mysql "easy-tiktok/db/mysql"
	"easy-tiktok/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.VideoServer
}

func (s Server) Feed(ctx context.Context, request *proto.DouyinFeedRequest) (*proto.DouyinFeedResponse, error) {
	select {
	//判断请求是否被取消
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request is canceled")
	default:
		// 继续执行
	}
	db := Mysql.GetDB()
	var list []Mysql.Video
	if db.Preload("UserMsg").Find(&list).Error != nil {
		return nil, nil
	}
	m := map[int64]bool{}
	likes := map[int64]bool{}
	if *request.Token != "" {
		userId := util.GetUserId(*request.Token)
		var follow []Mysql.Follow
		if db.Where("follower =?", userId).Find(&follow).Error == nil {
			for _, v := range follow {
				m[v.BeFollowed.Int64] = true
			}
		}
		var like []Mysql.Like
		if db.Where("liker_id =?", userId).Find(&like).Error == nil {
			for _, v := range like {
				likes[v.VideoID] = true
			}
		}
	}
	nextTime := int64(1)
	var list2 []*proto.Video
	for _, v := range list {
		if v.CreatedAt.Unix() > nextTime {
			nextTime = v.CreatedAt.Unix()
		}
		b := m[v.UserMsg.ID]
		b2 := likes[v.ID]
		list2 = append(list2, &proto.Video{
			Id: &v.ID,
			Author: &proto2.User{
				Id:              &v.UserMsg.ID,
				Name:            &v.UserMsg.Username,
				FollowCount:     &v.UserMsg.FollowCount,
				FollowerCount:   &v.UserMsg.FollowerCount,
				IsFollow:        &b,
				Avatar:          &v.UserMsg.Avatar.String,
				BackgroundImage: &v.UserMsg.BackgroundImage.String,
				Signature:       &v.UserMsg.Signature.String,
				TotalFavorited:  &v.UserMsg.TotalFavorited.Int64,
				WorkCount:       &v.UserMsg.WorkCount,
				FavoriteCount:   &v.UserMsg.FavoriteCount,
			},
			PlayUrl:       &v.PlayUrl,
			CoverUrl:      &v.CoverUrl,
			FavoriteCount: &v.FavoriteCount,
			CommentCount:  &v.CommentCount,
			IsFavorite:    &b2,
			Title:         &v.Title,
		})
	}
	return &proto.DouyinFeedResponse{
		StatusCode: &Mysql.S.Ok,
		StatusMsg:  &Mysql.S.OkMsg,
		VideoList:  list2,
		NextTime:   &nextTime,
	}, nil
}
