package model

import orm "easy-tiktok/db/mysql"

// Message //
// 消息model
// Author lql
type Message struct {
	orm.Model
	FromUserID int64
	ToUserId   int64
	Content    string
}

func (Message) TableName() string {
	return MESSAGE_TABLE
}

const MESSAGE_TABLE = "message"
