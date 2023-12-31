package main

import (
	"easy-tiktok/apps/app/config"
	"easy-tiktok/apps/interaction/internal"
	"easy-tiktok/apps/interaction/proto"
	"easy-tiktok/util/etcd"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	logger := logrus.New()
	config.Initial()
	ctx := context.Background()
	listen, err := net.Listen(config.C.NetworkType, config.C.InteractionHost)
	if err != nil {
		println("连接interaction服务不是很成功")
		println(err)
		return
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(grpclogrus.UnaryServerInterceptor(logrus.NewEntry(logger))),
		grpc.StreamInterceptor(grpclogrus.StreamServerInterceptor(logrus.NewEntry(logger))))
	proto.RegisterInteractionServer(server, internal.Server{})
	if err := etcd.Register(ctx, "interaction", config.C.InteractionHost); err != nil {
		println("Ectd注册服务失败")
		return
	}
	// TODO 应该在生产环境中禁用 reflection 功能，并且只在开发和测试阶段使用.
	reflection.Register(server)
	if err := server.Serve(listen); err != nil {
		return
	}
	defer server.Stop()
}
