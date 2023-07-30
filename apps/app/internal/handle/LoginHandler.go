package handle

import (
	"easy-tiktok/apps/app/internal/rpc"
	"easy-tiktok/apps/user/user"
	"github.com/gin-gonic/gin"
)

func loginHandler(c *gin.Context) {
	var req user.DouyinUserLoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, err)
		return
	}
	userRpc := rpc.GetUserRpc()
	userRpc.Login(rpc.GetCtx(), &req)
	//client.Login()
}
