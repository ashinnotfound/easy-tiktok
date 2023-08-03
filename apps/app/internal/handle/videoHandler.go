package handle

import (
	context2 "context"
	"easy-tiktok/apps/app/internal/rpc"
	video "easy-tiktok/apps/video/proto"
	"github.com/gin-gonic/gin"
	"io"
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
		context.JSON(400, err)
		return
	}
	context.JSON(200, &feed)
}

func videoActionHandler(context *gin.Context) {

	title := context.PostForm("title")
	token := context.PostForm("token")

	formFile, _, err2 := context.Request.FormFile("data")
	if err2 != nil {
		panic(err2)
	}

	videoFile, err3 := io.ReadAll(formFile)

	defer formFile.Close()
	if err3 != nil {
		panic(err3)
	}

	req := video.DouyinPublishActionRequest{
		Token: &token,
		Data:  videoFile,
		Title: &title,
	}
	videoRpc := rpc.GetVideoRpc()
	action, err := videoRpc.Action(context, &req)
	if err != nil {
		panic(err)
		context.JSON(400, err)
		return
	}
	context.JSON(200, &action)
}

func videoListHandler(context *gin.Context) {

	token := context.Query("token")
	userId, err := strconv.ParseInt(context.Query("user_id"), 10, 64)
	if err != nil {
		panic(err)
	}

	req := video.DouyinPublishListRequest{
		UserId: &userId,
		Token:  &token,
	}

	videoRpc := rpc.GetVideoRpc()
	list, err2 := videoRpc.List(context, &req)
	if err2 != nil {
		panic(err2)
		context.JSON(400, err)
		return
	}

	context.JSON(200, &list)

}
