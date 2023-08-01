package Mysql

import (
	"database/sql"
)

type UserEntry struct {
	Model
	Username string `gorm:"column:username;default:;NOT NULL;comment:'用户名'"`
	Password string `gorm:"column:password;default:;NOT NULL;comment:'用户密码，MD5加密'"`
}

type UserMsg struct {
	Model
	FollowCount     int64          `gorm:"column:follow_count;default:0;NOT NULL;comment:'关注数'"`
	FollowerCount   int64          `gorm:"column:follower_count;default:0;NOT NULL;comment:'粉丝数'"`
	Avatar          sql.NullString `gorm:"column:avatar;comment:'头像'"`
	BackgroundImage sql.NullString `gorm:"column:background_image;comment:'背景图片'"`
	Signature       sql.NullString `gorm:"column:signature;comment:'个性签名'"`
	TotalFavorited  sql.NullInt64  `gorm:"column:total_favorited;comment:'获得赞数'"`
	WorkCount       int64          `gorm:"column:work_count;default:0;NOT NULL;comment:'作品数'"`
	FavoriteCount   int64          `gorm:"column:favorite_count;default:0;NOT NULL;comment:'喜欢数'"`
	Username        string         `gorm:"column:username;NOT NULL"`
}

type Follow struct {
	Model
	BeFollowed sql.NullInt64 `gorm:"column:be_followed"`
	Follower   sql.NullInt64 `gorm:"column:follower"`
}

func (f *Follow) TableName() string {
	return "follow"
}

func (u *UserMsg) TableName() string {
	return "user_msg"
}

func (u *UserEntry) TableName() string {
	return "user_entry"
}
