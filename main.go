package main

import Mysql "easy-tiktok/db/mysql"

func main() {

	db := Mysql.GetDB()
	var list []Mysql.Video
	if db.Preload("UserMsg").Find(&list).Error != nil {
		return
	}
	println(list[0].UserMsg.ID)

}
