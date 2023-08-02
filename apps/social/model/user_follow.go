package model

import orm "easy-tiktok/db/mysql"

type UserFollow struct {
	orm.Model
	UserId   int64 `gorm:"column:user_id;NOT NULL;comment:'用户id'"`
	FollowId int64 `gorm:"column:follow_id;NOT NULL;comment:'关注用户id'"`
	Status   int32 `gorm:"column:status;NOT NULL;comment:'关注状态'"`
}

const USER_FOLLOW_TABLE = "user_follow"
