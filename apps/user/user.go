package main

import (
	"easy-tiktok/apps/user/internal/logic"
	"easy-tiktok/apps/user/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return
	}
	server := grpc.NewServer()
	user.RegisterUserServer(server, logic.Server{})
	reflection.Register(server)
	if err := server.Serve(listen); err != nil {
		return
	}
}
