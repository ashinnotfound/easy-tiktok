package rpc

import (
	"easy-tiktok/apps/user/user"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

var ctx context.Context
var userRpc user.UserClient

func init() {
	dial, err := grpc.Dial("localhost:50051", grpc.WithBlock())
	if err != nil {
		return
	}
	// 延迟关闭连接
	defer dial.Close()
	var cancel context.CancelFunc
	// 初始化上下文，设置请求超时时间为1秒
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	// 延迟关闭请求会话
	defer cancel()
	// 初始化Rpc服务客户端
	userRpc = user.NewUserClient(dial)
	//  videoRpc = video.NewVideoClient(dial)
}

func GetCtx() context.Context {
	return ctx
}

func GetUserRpc() user.UserClient {
	return userRpc
}
