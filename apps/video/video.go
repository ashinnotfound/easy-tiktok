package main

import (
	"easy-tiktok/apps/video/internal/logic"
	"easy-tiktok/apps/video/proto"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	logger := logrus.New()
	listen, err := net.Listen("tcp", "localhost:8091")
	if err != nil {
		println("连接不是很成功")
		println(err)
		return
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger))),
		grpc.StreamInterceptor(grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logger))))
	proto.RegisterVideoServer(server, logic.Server{})
	reflection.Register(server)
	if err := server.Serve(listen); err != nil {
		return
	}
	defer server.Stop()
}
