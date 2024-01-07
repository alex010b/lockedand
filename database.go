package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func dbConnect() {
	dsn := "root:ImABall!@#@@(@11212@tcp(127.0.0.1:3306)/lockedand?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	fmt.Println("connected to db")

	// AutoMigrate will create the table if it does not exist
	db.AutoMigrate(&Person{})
}
