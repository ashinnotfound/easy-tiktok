package rpc

import (
	"easy-tiktok/apps/app/config"
	interaction "easy-tiktok/apps/interaction/proto"
	user "easy-tiktok/apps/user/proto"
	video "easy-tiktok/apps/video/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userRpc user.UserClient
var videoRpc video.VideoClient
var interactionRpc interaction.InteractionClient

func Initial() {
	dial, err := grpc.Dial(config.C.UserHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	dial1, err := grpc.Dial(config.C.VideoHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	interactionDial, err := grpc.Dial(config.C.InteractionHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	// 初始化Rpc服务客户端
	userRpc = user.NewUserClient(dial)
	videoRpc = video.NewVideoClient(dial1)
	interactionRpc = interaction.NewInteractionClient(interactionDial)
}

func GetUserRpc() user.UserClient {
	return userRpc
}

func GetVideoRpc() video.VideoClient {
	return videoRpc
}

func GetInteractionRpc() interaction.InteractionClient {
	return interactionRpc
}
