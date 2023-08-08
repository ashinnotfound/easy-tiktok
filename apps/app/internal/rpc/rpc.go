package rpc

import (
	"easy-tiktok/apps/app/config"
	interaction "easy-tiktok/apps/interaction/proto"
	social "easy-tiktok/apps/social/proto"
	user "easy-tiktok/apps/user/proto"
	video "easy-tiktok/apps/video/proto"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userRpc user.UserClient
var videoRpc video.VideoClient
var interactionRpc interaction.InteractionClient
var socialRpc social.SocialClient

func Initial() {
	cli, err := clientv3.NewFromURL(config.C.Etcd)
	builder, err := resolver.NewBuilder(cli)
	userDial, err := grpc.Dial(config.C.UserService, grpc.WithResolvers(builder), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	videoDial, err := grpc.Dial(config.C.VideoService, grpc.WithResolvers(builder), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	interactionDial, err := grpc.Dial(config.C.InteractionService, grpc.WithResolvers(builder), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	socialDial, err := grpc.Dial(config.C.SocialService, grpc.WithResolvers(builder), grpc.WithTransportCredentials(insecure.NewCredentials()))
	// 初始化Rpc服务客户端
	userRpc = user.NewUserClient(userDial)
	videoRpc = video.NewVideoClient(videoDial)
	interactionRpc = interaction.NewInteractionClient(interactionDial)
	socialRpc = social.NewSocialClient(socialDial)
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

// GetSocialRpc //
// 获取社交模块的rpc客户端
func GetSocialRpc() social.SocialClient {
	return socialRpc
}
