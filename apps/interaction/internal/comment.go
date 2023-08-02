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
	"time"
)

// Comment POST /douyin/comment/action/ - 评论操作  登录用户对视频进行评论。
func (Server) Comment(ctx context.Context, request *proto.DouyinCommentActionRequest) (*proto.DouyinCommentActionResponse, error) {
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
	// 查找视频
	var video Mysql.Video
	if err := db.First(&video, request.GetVideoId()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.InvalidArgument, "Comment::invalid videoId")
		} else {
			return nil, status.Error(codes.Aborted, "Comment::database exception")
		}
	}
	// 待操作数
	var numToAdd int64
	if request.GetActionType() == 1 {
		numToAdd = 1
	} else if request.GetActionType() == 2 {
		numToAdd = -1
	} else {
		return nil, status.Error(codes.InvalidArgument, "Comment::invalid actionType")
	}
	// 开启事务
	tx := db.Begin()
	// 视频评论数+-1
	if tx.Model(&video).Update("comment_count", video.CommentCount+numToAdd).Error == nil {
		// 更新评论内容
		if request.GetActionType() == 1 {
			// 查找用户信息
			var userMsg Mysql.UserMsg
			if err := tx.First(&userMsg, userId).Error; err != nil {
				// 查找不到说明当前用户不存在
				if errors.Is(err, gorm.ErrRecordNotFound) {
					tx.Rollback()
					return nil, status.Error(codes.Unauthenticated, "Comment::invalid tokenId")
				} else {
					tx.Rollback()
					return nil, status.Error(codes.Aborted, "Comment::database exception")
				}
			}
			// 创建评论记录
			comment := Mysql.Comment{
				VideoID:    request.GetVideoId(),
				Content:    request.GetCommentText(),
				CreateDate: time.Now().Format("01-02"),
				UserMsg:    userMsg,
			}
			if tx.Create(&comment).Error != nil {
				tx.Rollback()
				return nil, status.Error(codes.Aborted, "Comment::database exception")
			}
			// 提交事务
			tx.Commit()
			statusMsg := "评论成功!!!"
			isFollow := false
			return &proto.DouyinCommentActionResponse{
				StatusCode: &Mysql.S.Ok,
				StatusMsg:  &statusMsg,
				Comment: &proto.Comment{
					Id: &comment.ID,
					User: &proto.User{
						Id:              &userMsg.ID,
						Name:            &userMsg.Username,
						FollowCount:     &userMsg.FollowCount,
						FollowerCount:   &userMsg.FollowerCount,
						IsFollow:        &isFollow,
						Avatar:          &userMsg.Avatar.String,
						BackgroundImage: &userMsg.BackgroundImage.String,
						Signature:       &userMsg.Signature.String,
						TotalFavorited:  &userMsg.TotalFavorited.Int64,
						WorkCount:       &userMsg.WorkCount,
						FavoriteCount:   &userMsg.FavoriteCount,
					},
					Content:    request.CommentText,
					CreateDate: &comment.CreateDate,
				},
			}, nil
		} else {
			// 删除评论记录
			if tx.Delete(&Mysql.Comment{}, request.GetCommentId()).Error != nil {
				tx.Rollback()
				return nil, status.Error(codes.Aborted, "Comment::database exception")
			}
			// 提交事务
			tx.Commit()
			statusMsg := "删除评论成功!!!"
			return &proto.DouyinCommentActionResponse{
				StatusCode: &Mysql.S.Ok,
				StatusMsg:  &statusMsg,
			}, nil
		}
	}
	// 业务失败
	tx.Rollback()
	return nil, status.Error(codes.Aborted, "Comment::operation failed")
}

// GetCommentList GET /douyin/comment/list/ - 视频评论列表  查看视频的所有评论，按发布时间倒序。
func (Server) GetCommentList(ctx context.Context, request *proto.DouyinCommentListRequest) (*proto.DouyinCommentListResponse, error) {
	select {
	// 判断请求是否被取消
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "Favorite::request is canceled")
	default:
		// 继续执行
	}

	// 查找数据库中当前视频的所有评论
	var comment []Mysql.Comment
	db := Mysql.GetDB()
	if db.Where("video_id = ?", request.GetVideoId()).Find(&comment).Error != nil {
		return nil, status.Error(codes.Aborted, "GetCommentList::database exception")
	}

	// 填充返回值评论列表
	var commentList []*proto.Comment
	for _, v := range comment {
		// 评论用户信息
		userMsg := v.UserMsg
		// 查询当前用户是否关注了评论用户
		var follow Mysql.Follow
		isFollow := true
		// 验证token
		if request.GetToken() == "" {
			// 没登陆则没关注
			isFollow = false
		} else {
			userId := util.GetUserId(request.GetToken())
			if err := db.Where("be_followed = ? AND follower = ?", userMsg.ID, userId).First(&follow).Error; err != nil {
				// 找不到记录说明没关注
				if errors.Is(err, gorm.ErrRecordNotFound) {
					isFollow = false
				} else {
					return nil, status.Error(codes.Aborted, "GetCommentList::database exception")
				}
			}
		}
		commentList = append(commentList, &proto.Comment{
			Id: &v.ID,
			User: &proto.User{
				Id:              &userMsg.ID,
				Name:            &userMsg.Username,
				FollowCount:     &userMsg.FollowCount,
				FollowerCount:   &userMsg.FollowerCount,
				IsFollow:        &isFollow,
				Avatar:          &userMsg.Avatar.String,
				BackgroundImage: &userMsg.BackgroundImage.String,
				Signature:       &userMsg.Signature.String,
				TotalFavorited:  &userMsg.TotalFavorited.Int64,
				WorkCount:       &userMsg.WorkCount,
				FavoriteCount:   &userMsg.FavoriteCount,
			},
			Content:    &v.Content,
			CreateDate: &v.CreateDate,
		})
	}
	statusMsg := "获取视频评论列表成功"
	return &proto.DouyinCommentListResponse{
		StatusCode:  &Mysql.S.Ok,
		StatusMsg:   &statusMsg,
		CommentList: commentList,
	}, nil
}
