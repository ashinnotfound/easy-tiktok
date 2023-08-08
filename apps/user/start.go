package main

import (
	"easy-tiktok/apps/app/config"
	"easy-tiktok/apps/user/internal/logic"
	"easy-tiktok/apps/user/proto"
	"easy-tiktok/util/etcd"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
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
	listen, err := net.Listen("tcp", config.C.UserHost)
	if err != nil {
		println("连接不是很成功")
		println(err)
		return
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger))),
		grpc.StreamInterceptor(grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logger))))
	proto.RegisterUserServer(server, logic.Server{})
	if err := etcd.Register(ctx, "user", config.C.UserHost); err != nil {
		println("Ectd注册服务失败")
		return
	}
	reflection.Register(server)
	if err := server.Serve(listen); err != nil {
		return
	}
	defer server.Stop()

}
