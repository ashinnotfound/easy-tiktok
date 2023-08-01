package handle

import (
	context2 "context"
	"easy-tiktok/apps/app/internal/rpc"
	user "easy-tiktok/apps/user/proto"
	"github.com/gin-gonic/gin"
	"strconv"
)

func loginHandler(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	req := user.DouyinUserLoginRequest{
		Username: &username,
		Password: &password,
	}
	userRpc := rpc.GetUserRpc()
	login, err := userRpc.Login(context2.Background(), &req)
	if err != nil {
		c.JSON(400, login)
		return
	}
	c.JSON(200, login)
}

func userInfoHandler(context *gin.Context) {
	userId := context.Query("user_id")
	parseInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return
	}
	token := context.Query("token")
	req := user.DouyinUserRequest{
		UserId: &parseInt,
		Token:  &token,
	}
	userRpc := rpc.GetUserRpc()
	userInfo, err := userRpc.GetUserInfo(context2.Background(), &req)
	if err != nil {
		context.JSON(400, err)
		return
	}
	context.JSON(200, userInfo)
}

func registerHandler(context *gin.Context) {
	username := context.Query("username")
	password := context.Query("password")
	req := user.DouyinUserRegisterRequest{
		Username: &username,
		Password: &password,
	}
	userRpc := rpc.GetUserRpc()
	register, err := userRpc.Register(context2.Background(), &req)
	if err != nil {
		context.JSON(400, register)
		return
	}
	context.JSON(200, register)
}
