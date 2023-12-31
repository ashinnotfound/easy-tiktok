package logic

import (
	"database/sql"
	"easy-tiktok/apps/social/model"
	"easy-tiktok/apps/user/proto"
	Mysql "easy-tiktok/db/mysql"
	"easy-tiktok/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Server struct {
	proto.UserServer
}

var wrong int64 = 1

func (l Server) Login(ctx context.Context, request *proto.DouyinUserLoginRequest) (*proto.DouyinUserLoginResponse, error) {
	select {
	//判断请求是否被取消
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request is canceled")
	default:
		// 继续执行
	}
	model := Mysql.GetDB()
	user := Mysql.UserEntry{}
	tx := model.Where("username = ? and password = ?", request.GetUsername(), request.GetPassword()).First(&user)
	if tx.Error != nil {
		str := "用户名或密码错误"
		return &proto.DouyinUserLoginResponse{
			StatusCode: &Mysql.S.Bad,
			StatusMsg:  &str,
			UserId:     &wrong,
			Token:      &str,
		}, nil
	}
	token, err := util.GetToken(user.ID)
	if err != nil {
		return nil, err
	}
	return &proto.DouyinUserLoginResponse{
		StatusMsg:  &Mysql.S.OkMsg,
		StatusCode: &Mysql.S.Ok,
		UserId:     &user.ID,
		Token:      &token,
	}, nil
}

func (l Server) Register(ctx context.Context, request *proto.DouyinUserRegisterRequest) (*proto.DouyinUserRegisterResponse, error) {
	select {
	//判断请求是否被取消
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request is canceled")
	default:
		// 继续执行
	}
	model := Mysql.GetDB()
	var entry Mysql.UserEntry
	tx := model.Where("username = ?", request.GetUsername()).First(&entry)
	if tx.Error == nil || tx.RowsAffected != 0 {
		str := "用户名重复"
		return &proto.DouyinUserRegisterResponse{
			StatusCode: &Mysql.S.Bad,
			StatusMsg:  &str,
			UserId:     &wrong,
			Token:      &str,
		}, nil
	}
	entry = Mysql.UserEntry{
		Username: request.GetUsername(),
		Password: request.GetPassword(),
	}
	err := model.Transaction(func(tx *gorm.DB) error {
		if err := model.Create(&entry).Error; err != nil {
			return nil
		}
		msg := userMsgInit(entry.Username)
		if err := model.Create(&msg).Error; err != nil {
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	token, err := util.GetToken(entry.ID)
	if err != nil {
		return nil, err
	}
	return &proto.DouyinUserRegisterResponse{
		StatusCode: &Mysql.S.Ok,
		StatusMsg:  &Mysql.S.OkMsg,
		UserId:     &entry.ID,
		Token:      &token,
	}, nil
}

func (l Server) GetUserInfo(ctx context.Context, request *proto.DouyinUserRequest) (*proto.DouyinUserResponse, error) {
	select {
	//判断请求是否被取消
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request is canceled")
	default:
		// 继续执行
	}
	id := request.GetUserId()
	viewerId := util.GetUserId(request.GetToken())
	db := Mysql.GetDB()
	user := Mysql.UserMsg{}
	var follow model.UserFollow
	db.Where("id = ?", id).First(&user)
	var result bool
	if db.Where("follow_id = ? AND user_id = ?", id, viewerId).First(&follow).Error != nil {
		result = false
	} else {
		result = true
	}
	return &proto.DouyinUserResponse{
		StatusCode: &Mysql.S.Ok,
		StatusMsg:  &Mysql.S.OkMsg,
		User: &proto.User{
			Id:              &user.ID,
			Name:            &user.Username,
			FollowCount:     &user.FollowCount,
			FollowerCount:   &user.FollowerCount,
			IsFollow:        &result,
			Avatar:          &user.Avatar.String,
			BackgroundImage: &user.BackgroundImage.String,
			Signature:       &user.Signature.String,
			TotalFavorited:  &user.TotalFavorited.Int64,
			WorkCount:       &user.WorkCount,
			FavoriteCount:   &user.FavoriteCount,
		},
	}, nil
}

func userMsgInit(username string) Mysql.UserMsg {
	return Mysql.UserMsg{
		FollowCount:     0,
		FollowerCount:   0,
		Avatar:          sql.NullString{String: "https:ipfs.io/ipfs/bafkreiacrj7wlkvtbckd3cemrkcl3tu73upwiacu5debjjn6viyepaghka", Valid: true},
		BackgroundImage: sql.NullString{String: "https:ipfs.io/ipfs/bafkreiacrj7wlkvtbckd3cemrkcl3tu73upwiacu5debjjn6viyepaghka", Valid: true},
		Signature:       sql.NullString{String: "我想重新做人", Valid: true},
		TotalFavorited:  sql.NullInt64{Valid: true},
		WorkCount:       0,
		FavoriteCount:   0,
		Username:        username,
	}
}
