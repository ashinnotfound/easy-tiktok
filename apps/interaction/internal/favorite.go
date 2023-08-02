package internal

import (
	"context"
	"easy-tiktok/apps/interaction/proto"
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

	// 开始事务
	tx := Mysql.GetDB().Begin()
	// 查找当前用户
	var userMsg Mysql.UserMsg
	tx.First(&userMsg, userId)
	// 查找视频
	var video Mysql.Video
	tx.First(&video, request.GetVideoId())
	// 查找视频发布者
	var videoMakerMsg Mysql.UserMsg
	tx.First(&videoMakerMsg, video.UserMsgID)

	// 查找视频点赞记录 判断点赞/取消点赞
	var like Mysql.Like
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
			if tx.Model(&videoMakerMsg).Update("total_favorited", videoMakerMsg.TotalFavorited.Int64+numToAdd).Error == nil {
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
