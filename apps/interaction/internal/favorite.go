package internal

import (
	"context"
	"easy-tiktok/apps/constant"
	"easy-tiktok/apps/interaction/proto"
	"easy-tiktok/apps/social/model"
	Mysql "easy-tiktok/db/mysql"
	"easy-tiktok/util"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Server struct {
	proto.InteractionServer
}

// Favorite POST /douyin/favorite/action/ - 赞操作  登录用户对视频的点赞和取消点赞操作。
func (Server) Favorite(ctx context.Context, request *proto.DouyinFavoriteActionRequest) (*proto.DouyinFavoriteActionResponse, error) {
	select {
	// 判断请求是否被取消
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "Favorite::request is canceled")
	default:
		// 继续执行
	}

	// 验证token
	if request.GetToken() == "" {
		return nil, status.Error(codes.Unauthenticated, "Favorite::invalid token")
	}
	userId := util.GetUserId(request.GetToken())

	db := Mysql.GetDB()
	// 查找当前用户
	var userMsg Mysql.UserMsg
	if err := db.First(&userMsg, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.InvalidArgument, "Favorite::invalid token")
		} else {
			return nil, status.Error(codes.Aborted, "Favorite::database exception")
		}
	}
	// 查找视频和视频发布者
	var video Mysql.Video
	if err := db.Preload("UserMsg").First(&video, request.GetVideoId()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.InvalidArgument, "Favorite::invalid videoId")
		} else {
			return nil, status.Error(codes.Aborted, "Favorite::database exception")
		}
	}

	// 查找视频点赞记录 判断点赞/取消点赞
	var like Mysql.Like
	// 开始事务
	tx := db.Begin()
	if err := tx.Where("video_id = ? AND liker_id = ?", request.GetVideoId(), userId).First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 进行点赞操作
			if request.GetActionType() != 1 {
				return nil, status.Error(codes.Aborted, "Favorite::invalid actionType")
			}
			// 创建点赞记录
			if tx.Create(&Mysql.Like{VideoID: request.GetVideoId(), LikerID: userId}).Error != nil {
				// 业务失败
				tx.Rollback()
				return nil, status.Error(codes.Aborted, "Favorite::database exception")
			}
		} else {
			return nil, status.Error(codes.Aborted, "Favorite::database exception")
		}
	} else {
		// 进行取消点赞操作
		if request.GetActionType() != 2 {
			return nil, status.Error(codes.Aborted, "Favorite::invalid actionType")
		}
		// 删除点赞记录
		if tx.Delete(&like).Error != nil {
			// 业务失败
			tx.Rollback()
			return nil, status.Error(codes.Aborted, "Favorite::database exception")
		}
	}
	// 待操作数
	var numToAdd int64
	if request.GetActionType() == 1 {
		numToAdd = 1
	} else if request.GetActionType() == 2 {
		numToAdd = -1
	} else {
		return nil, status.Error(codes.InvalidArgument, "Favorite::invalid actionType")
	}
	// 用户点赞数+-1
	if tx.Model(&userMsg).Update("favorite_count", userMsg.FavoriteCount+numToAdd).Error == nil {
		// 视频的点赞总数+-1
		if tx.Model(&video).Update("favorite_count", video.FavoriteCount+numToAdd).Error == nil {
			// 视频发布者获得赞数+-1
			if tx.Model(&video.UserMsg).Update("total_favorited", video.UserMsg.TotalFavorited.Int64+numToAdd).Error == nil {
				// 提交事务
				tx.Commit()
				return &proto.DouyinFavoriteActionResponse{
					StatusCode: &Mysql.S.Ok,
					StatusMsg:  &Mysql.S.OkMsg,
				}, nil
			}
		}
	}
	// 业务失败
	tx.Rollback()
	return nil, status.Error(codes.Aborted, "Favorite::operation failed")
}

// GetFavoriteList GET /douyin/favorite/list/ - 喜欢列表  登录用户的所有点赞视频。
func (Server) GetFavoriteList(ctx context.Context, request *proto.DouyinFavoriteListRequest) (*proto.DouyinFavoriteListResponse, error) {
	select {
	// 判断请求是否被取消
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "Favorite::request is canceled")
	default:
		// 继续执行
	}

	//// 验证token
	//if request.GetToken() == "" {
	//	return nil, status.Error(codes.Unauthenticated, "GetFavoriteList::invalid token")
	//}
	//userId := util.GetUserId(request.GetToken())
	//if userId != request.GetUserId() {
	//	return nil, status.Error(codes.Unauthenticated, "GetFavoriteList::invalid token")
	//}
	userId := request.GetUserId()

	db := Mysql.GetDB()
	// 查找当前用户的点赞视频
	var like []Mysql.Like
	if err := db.Where("liker_id = ?", userId).Find(&like).Error; err != nil {
		return nil, status.Error(codes.Aborted, "GetFavoriteList::database exception")
	}
	// 查找数据库的视频列表
	var video []Mysql.Video
	for _, v := range like {
		var currentVideo Mysql.Video
		if err := db.First(&currentVideo, v.VideoID).Error; err != nil {
			// TODO 后续对查找失败的视频做进一步处理
			continue
		}
		video = append(video, currentVideo)
	}
	// 填充返回值视频列表
	var videoList []*proto.Video
	isFavorite := true
	for i := range video {
		v := video[i]
		var follow model.UserFollow
		isFollow := true
		if err := db.Where("follow_id = ? AND user_id = ? AND status = ?", v.UserMsgID, userId, constant.RELATION_FOLLOW).First(&follow).Error; err != nil {
			// 找不到记录说明没关注
			if errors.Is(err, gorm.ErrRecordNotFound) {
				isFollow = false
			} else {
				return nil, status.Error(codes.Aborted, "GetFavoriteList::database exception")
			}
		}

		videoList = append(videoList, &proto.Video{
			Id: &v.ID,
			Author: &proto.User{
				Id:              &v.UserMsg.ID,
				Name:            &v.UserMsg.Username,
				FollowCount:     &v.UserMsg.FollowCount,
				FollowerCount:   &v.UserMsg.FollowerCount,
				IsFollow:        &isFollow,
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
			IsFavorite:    &isFavorite,
			Title:         &v.Title,
		})
	}

	statusMsg := "获取用户喜欢视频列表成功"
	return &proto.DouyinFavoriteListResponse{
		StatusCode: &Mysql.S.Ok,
		StatusMsg:  &statusMsg,
		VideoList:  videoList,
	}, nil
}
