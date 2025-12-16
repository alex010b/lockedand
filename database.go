package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func dbConnect() {
	dsn := "root:"
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	fmt.Println("connected to db")

	// AutoMigrate will create the table if it does not exist
	db.AutoMigrate(&Person{})
}
