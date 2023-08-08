package main

import (
	"easy-tiktok/apps/app/config"
	"easy-tiktok/apps/global"
	"easy-tiktok/apps/initialize"
	"easy-tiktok/apps/social/internal/logic"
	pb "easy-tiktok/apps/social/proto"
	"easy-tiktok/util/etcd"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	initialize.LogInit()
	config.Initial()
	listen, err := net.Listen(config.C.NetworkType, config.C.SocialHost)
	ctx := context.Background()
	if err != nil {
		global.LOGGER.Errorf("连接soical服务不是很成功,reason: %v\n", err)
		return
	}
	if err != nil {
		global.LOGGER.Infof("user_follow table can not create")
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(grpclogrus.UnaryServerInterceptor(logrus.NewEntry(global.LOGGER))),
		grpc.StreamInterceptor(grpclogrus.StreamServerInterceptor(logrus.NewEntry(global.LOGGER))))
	pb.RegisterSocialServer(server, &logic.SocialServerImpl{})
	if err := etcd.Register(ctx, "social", config.C.SocialHost); err != nil {
		println("Ectd注册服务失败")
		return
	}
	// TODO 应该在生产环境中禁用 reflection 功能，并且只在开发和测试阶段使用.
	reflection.Register(server)
	global.LOGGER.Infof("social RPC server listen")
	if err := server.Serve(listen); err != nil {
		return
	}
	defer server.Stop()
}
