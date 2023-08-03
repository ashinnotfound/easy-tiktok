package main

import (
	Mysql "easy-tiktok/db/mysql"
	"fmt"
)

func main() {

	db := Mysql.GetDB()
	var list []Mysql.Video
	//db.Preload("UserMsg").Where("").Find(&list)
	db.Preload("UserMsg").Where("userMsg_id = ?", 20).Find(&list)

	for index, _ := range list {
		fmt.Println(list[index])
		fmt.Println("\n")
	}

}
