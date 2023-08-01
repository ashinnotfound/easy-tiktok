package handle

import (
	context2 "context"
	"easy-tiktok/apps/app/internal/rpc"
	video "easy-tiktok/apps/video/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func videoFeedHandler(context *gin.Context) {
	token := context.Query("token")
	time := context.Query("latest_time")
	parseInt, err := strconv.ParseInt(time, 10, 64)

	req := video.DouyinFeedRequest{
		LatestTime: &parseInt,
		Token:      &token,
	}
	videoRpc := rpc.GetVideoRpc()
	feed, err := videoRpc.Feed(context2.Background(), &req)
	if err != nil {
		panic(err)
		context.JSON(200, err)
		return
	}
	context.JSON(http.StatusOK, &feed)
}
