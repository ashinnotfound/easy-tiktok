package handle

import (
	"bytes"
	context2 "context"
	"easy-tiktok/apps/app/internal/rpc"
	video "easy-tiktok/apps/video/proto"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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

	//获取request以body形式传输的参数
	type actionBodyType struct {
		File  []byte `json:"data"`
		Token string `json:"token"`
		Title string `json:"title"`
	}
	var actionBody actionBodyType
	data, err := context.GetRawData()
	if err != nil {
		fmt.Println(err)
	}
	//把字节流重新放回body
	context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	//转成结构体类型进行参数的接收
	if err := json.NewDecoder(context.Request.Body).Decode(&actionBody); err != nil {
		panic(err)
		context.JSON(400, err)
	}
	req := video.DouyinPublishActionRequest{
		Token: &actionBody.Token,
		Data:  actionBody.File,
		Title: &actionBody.Title,
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

}
