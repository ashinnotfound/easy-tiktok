package main

import (
	"easy-tiktok/apps/app/config"
	"easy-tiktok/apps/app/internal/handle"
	"easy-tiktok/apps/app/internal/rpc"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	config.Initial()

	rpc.Initial()
	handle.Route(engine)
	err := engine.Run(":8888")
	if err != nil {
		return
	}
}
