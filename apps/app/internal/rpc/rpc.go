package rpc

import (
	"easy-tiktok/apps/app/internal/config"
	"easy-tiktok/apps/user/user"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

var userRpc user.UserClient

func init() {
	dial, err := grpc.Dial(config.C.RPC.Host, grpc.WithBlock())
	if err != nil {
		return
	}
	// 延迟关闭连接
	defer dial.Close()
	// 初始化Rpc服务客户端
	userRpc = user.NewUserClient(dial)
	//  videoRpc = video.NewVideoClient(dial)
}

func GetCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return ctx
}

func GetUserRpc() user.UserClient {
	return userRpc
}
