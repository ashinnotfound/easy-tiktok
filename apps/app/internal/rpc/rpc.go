package rpc

import (
	"easy-tiktok/apps/app/internal/config"
	user "easy-tiktok/apps/user/proto"
	video "easy-tiktok/apps/video/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userRpc user.UserClient
var videoRpc video.VideoClient

func Initial() {
	dial, err := grpc.Dial(config.C.UserHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	dial1, errr := grpc.Dial(config.C.VideoHost, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if errr != nil {
		panic(errr)
	}
	// 初始化Rpc服务客户端
	userRpc = user.NewUserClient(dial)
	videoRpc = video.NewVideoClient(dial1)
}

func GetUserRpc() user.UserClient {
	return userRpc
}

func GetVideoRpc() video.VideoClient {
	return videoRpc
}
