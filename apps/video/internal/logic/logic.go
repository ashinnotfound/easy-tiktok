package logic

import (
	"bytes"
	proto2 "easy-tiktok/apps/user/proto"
	"easy-tiktok/apps/video/proto"
	Mysql "easy-tiktok/db/mysql"
	"easy-tiktok/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	shell "github.com/ipfs/go-ipfs-api"
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
		userid := util.GetUserId(*request.Token)
		var follow []Mysql.Follow
		if db.Where("follower =?", userid).Find(&follow).Error == nil {
			for _, v := range follow {
				m[v.BeFollowed.Int64] = true
			}
		}
		var like []Mysql.Like
		if db.Where("liker_id =?", userid).Find(&like).Error == nil {
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

func (s Server) Action(ctx context.Context, request *proto.DouyinPublishActionRequest) (*proto.DouyinPublishActionResponse, error) {
	select {
	//判断请求是否被取消
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request is canceled")
	default:
		// 继续执行
	}

	//数据库对象
	model := Mysql.GetDB()

	//获取token
	token := request.GetToken()

	//视频标题
	title := request.GetTitle()

	//获取用户传输的视频
	videoPut := request.GetData()

	//IPFS得到cid
	sh := shell.NewShell("10.21.23.163:6666")
	videoHash, err := sh.Add(bytes.NewReader(videoPut))

	if err != nil {
		return nil, err
	}

	//得到userId
	userId := util.GetUserId(token)

	//封装结构体对象
	videoMsg := Mysql.Video{
		Model:         Mysql.Model{},
		UserMsgID:     userId,
		UserMsg:       Mysql.UserMsg{},
		PlayUrl:       "https://ipfs.ashinnotfound.top/ipfs/" + videoHash,
		CoverUrl:      "",
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}

	if err := model.Create(&videoMsg).Error; err != nil {
		return nil, model.Error
	}
	return &proto.DouyinPublishActionResponse{
		StatusCode: &Mysql.S.Ok,
		StatusMsg:  &Mysql.S.OkMsg,
	}, nil
}

func (s Server) List(ctx context.Context, request *proto.DouyinPublishListRequest) (*proto.DouyinPublishListResponse, error) {
	select {
	//判断请求是否被取消
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request is canceled")
	default:
		// 继续执行
	}

	//获取数据库对象
	model := Mysql.GetDB()

	//获取token
	token := request.GetToken()

	//通过token获取userid
	userId := util.GetUserId(token)

	var list []*proto.Video

	//根据userId在video表中查找
	model.Where("userMsg_id = ?", userId).Find(&list)
	return &proto.DouyinPublishListResponse{
		StatusCode: &Mysql.S.Ok,
		StatusMsg:  &Mysql.S.OkMsg,
		VideoList:  list,
	}, nil
}
