package handle

import (
	"easy-tiktok/apps/app/internal/rpc"
	interaction "easy-tiktok/apps/interaction/proto"
	"github.com/gin-gonic/gin"
	"strconv"
)

func favoriteHandler(context *gin.Context) {
	token := context.Query("token")
	videoId, _ := strconv.ParseInt(context.Query("video_id"), 10, 64)
	actionTypeInt64, _ := strconv.ParseInt(context.Query("action_type"), 10, 32)
	actionType := int32(actionTypeInt64)
	req := interaction.DouyinFavoriteActionRequest{
		Token:      &token,
		VideoId:    &videoId,
		ActionType: &actionType,
	}
	favorite, err := rpc.GetInteractionRpc().Favorite(context, &req)
	if err != nil {
		context.JSON(400, favorite)
		return
	}
	context.JSON(200, favorite)
}

func getFavoriteListHandler(context *gin.Context) {
	token := context.Query("token")
	userId, _ := strconv.ParseInt(context.Query("user_id"), 10, 64)

	req := interaction.DouyinFavoriteListRequest{
		Token:  &token,
		UserId: &userId,
	}

	favoriteList, err := rpc.GetInteractionRpc().GetFavoriteList(context, &req)
	if err != nil {
		context.JSON(400, favoriteList)
		return
	}
	context.JSON(200, favoriteList)
}

func commentHandler(context *gin.Context) {
	token := context.Query("token")
	videoId, _ := strconv.ParseInt(context.Query("video_id"), 10, 64)
	actionTypeInt64, _ := strconv.ParseInt(context.Query("action_type"), 10, 32)
	actionType := int32(actionTypeInt64)
	commentText := context.Query("comment_text")
	commentId, _ := strconv.ParseInt(context.Query("comment_id"), 10, 64)

	req := interaction.DouyinCommentActionRequest{
		Token:       &token,
		VideoId:     &videoId,
		ActionType:  &actionType,
		CommentText: &commentText,
		CommentId:   &commentId,
	}

	comment, err := rpc.GetInteractionRpc().Comment(context, &req)
	if err != nil {
		context.JSON(400, comment)
		return
	}
	context.JSON(200, comment)
}

func getCommentListHandler(context *gin.Context) {
	token := context.Query("token")
	videoId, _ := strconv.ParseInt(context.Query("video_id"), 10, 64)

	req := interaction.DouyinCommentListRequest{
		Token:   &token,
		VideoId: &videoId,
	}

	commentList, err := rpc.GetInteractionRpc().GetCommentList(context, &req)
	if err != nil {
		context.JSON(400, commentList)
		return
	}
	context.JSON(200, commentList)
}
