package Mysql

type Video struct {
	Model
	UserMsgID     int64 `gorm:"column:userMsg_id;NOT NULL;comment:'视频所属用户id'"`
	UserMsg       UserMsg
	PlayUrl       string `gorm:"column:play_url;NOT NULL;comment:'视频播放地址'"`
	CoverUrl      string `gorm:"column:cover_url;NOT NULL;comment:'视频封面地址'"`
	FavoriteCount int64  `gorm:"column:favorite_count;default:0;NOT NULL;comment:'视频的点赞总数'"`
	CommentCount  int64  `gorm:"column:comment_count;default:0;NOT NULL;comment:'视频的评论总数'"`
	Title         string `gorm:"column:title;NOT NULL;comment:'视频标题'"`
}
type Like struct {
	Model
	VideoID int64 `gorm:"column:video_id;NOT NULL;comment:'视频id'"`
	LikerID int64 `gorm:"column:liker_id;NOT NULL;comment:'点赞者id'"`
}

func (v *Video) TableName() string {
	return "video"
}
