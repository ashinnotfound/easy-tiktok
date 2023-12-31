package logic

import (
	"bytes"
	"easy-tiktok/apps/constant"
	"easy-tiktok/apps/social/model"
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
	var videos []Mysql.Video
	//预加载对应信息
	if db.Preload("UserMsg").Find(&videos).Error != nil {
		return nil, nil
	}
	followMap := map[int64]bool{}
	likesMap := map[int64]bool{}
	//判断用户是否登录
	//登录则显示对应关注点赞关系
	if *request.Token != "" {
		userid := util.GetUserId(*request.Token)
		var follow []model.UserFollow
		if db.Where("user_id = ? AND status = ?", userid, constant.RELATION_FOLLOW).Find(&follow).Error == nil {
			for _, v := range follow {
				followMap[v.FollowId] = true
			}
		}
		var like []Mysql.Like
		if db.Where("liker_id =?", userid).Find(&like).Error == nil {
			for _, v := range like {
				likesMap[v.VideoID] = true
			}
		}
	}
	//调用loadVideos函数
	videoList, nextTime := loadVideos(followMap, likesMap, videos)
	return &proto.DouyinFeedResponse{
		StatusCode: &Mysql.S.Ok,
		StatusMsg:  &Mysql.S.OkMsg,
		VideoList:  videoList,
		NextTime:   &nextTime,
	}, nil
}

func (s Server) Action(ctx context.Context, request *proto.DouyinPublishActionRequest) (*proto.DouyinPublishActionResponse, error) {

	//数据库对象,开启事务
	model := Mysql.GetDB().Begin()

	select {
	//判断请求是否被取消
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request is canceled")
	default:
		// 继续执行
	}
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
		model.Rollback()
		return nil, model.Error
	}

	var user Mysql.UserMsg
	//通过userid获取其对应的作品数,并更新(获得了videoHash证明视频上传成功,故可以执行更新作品数操作)
	model.Where("id=?", userId).Find(&user)
	//让发布者的作品数加一
	user.WorkCount += 1
	//保存修改
	if err := model.Table("user_msg").Where("id=?", userId).Update("work_count", user.WorkCount).Error; err != nil {
		model.Rollback()
	}

	//提交事务
	model.Commit()

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

	db := Mysql.GetDB()
	var videos []Mysql.Video
	if db.Preload("UserMsg").Where("userMsg_id=?", request.GetUserId()).Find(&videos).Error != nil {
		return nil, nil
	}
	followMap := map[int64]bool{}
	likesMap := map[int64]bool{}
	//判断用户是否登录
	//登录则显示对应关注点赞关系
	if *request.Token != "" {
		userid := util.GetUserId(*request.Token)
		var follow []model.UserFollow
		if db.Where("user_id = ? AND status = ?", userid, constant.RELATION_FOLLOW).Find(&follow).Error == nil {
			for _, v := range follow {
				followMap[v.FollowId] = true
			}
		}
		var like []Mysql.Like
		if db.Where("liker_id =?", userid).Find(&like).Error == nil {
			for _, v := range like {
				likesMap[v.VideoID] = true
			}
		}
	}

	//调用loadVideos函数
	videoList, _ := loadVideos(followMap, likesMap, videos)
	return &proto.DouyinPublishListResponse{
		StatusCode: &Mysql.S.Ok,
		StatusMsg:  &Mysql.S.OkMsg,
		VideoList:  videoList,
	}, nil
}

// 装载video

func loadVideos(followMap map[int64]bool, likesMap map[int64]bool, videos []Mysql.Video) ([]*proto.Video, int64) {
	var nextTime int64
	var videoList []*proto.Video
	for i, _ := range videos {
		{
			v := videos[i]
			if v.CreatedAt.Unix() > nextTime {
				nextTime = v.CreatedAt.Unix()
			}
			b := followMap[v.UserMsg.ID]
			b2 := likesMap[v.ID]
			videoList = append(videoList, &proto.Video{
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
	}
	return videoList, nextTime
}
