package main

import (
	"easy-tiktok/apps/app/internal/handle"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	handle.Route(engine)
	//engine.Use(middleware)
	engine.Run("localhost:8080")
}
