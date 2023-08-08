package main

import (
	"context"
	"easy-tiktok/apps/app/config"
	"easy-tiktok/apps/video/internal/logic"
	"easy-tiktok/apps/video/proto"
	"easy-tiktok/util/etcd"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	logger := logrus.New()
	ctx := context.Background()
	listen, err := net.Listen("tcp", config.C.VideoHost)
	if err != nil {
		println("连接不是很成功")
		println(err)
		return
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger))),
		grpc.StreamInterceptor(grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logger))), grpc.MaxRecvMsgSize(200*1024*1024))
	proto.RegisterVideoServer(server, logic.Server{})
	if err := etcd.Register(ctx, "video", config.C.VideoHost); err != nil {
		println("Ectd注册服务失败")
		return
	}
	reflection.Register(server)
	if err := server.Serve(listen); err != nil {
		return
	}
	defer server.Stop()
}
