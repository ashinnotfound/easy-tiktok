package logic

import (
	"easy-tiktok/apps/user/user"
	"golang.org/x/net/context"
)

type Server struct {
	user.UserServer
}

func (l Server) Login(ctx context.Context, request *user.DouyinUserLoginRequest) (*user.DouyinUserLoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l Server) Register(ctx context.Context, request *user.DouyinUserRegisterRequest) (*user.DouyinUserRegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (l Server) mustEmbedUnimplementedUserServer() {
	//TODO implement me
	panic("implement me")
}
