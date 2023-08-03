package Mysql

import (
	"easy-tiktok/apps/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var S = Status{}
var _db *gorm.DB

type Model struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Status struct {
	Ok    int32
	Bad   int32
	OkMsg string
}

func init() {
	var err error
	_db, err = gorm.Open(mysql.Open("ljh:password@tcp(10.21.23.42:3306)/douyin?charset=utf8&parseTime=true"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	sqlDB, _ := _db.DB()

	//设置数据库连接池参数
	sqlDB.SetMaxOpenConns(100) //设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
	S = Status{
		Ok:    0,
		Bad:   1,
		OkMsg: "",
	}
	global.DB = _db
}

func GetDB() *gorm.DB {
	return _db
}
