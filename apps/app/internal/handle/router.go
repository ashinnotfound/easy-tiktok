package handle

import "github.com/gin-gonic/gin"

func Route(e *gin.Engine) {
	e.POST("/douyin/user/login", loginHandler)
}
