package main

import (
	"easy-tiktok/apps/app/internal/config"
	"easy-tiktok/apps/app/internal/handle"
	"easy-tiktok/apps/app/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	handle.Route(engine)
	engine.Use(middleware.Jwt())
	engine.Run(config.C.Host)
}
