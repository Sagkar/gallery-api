package db

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDBConnection() *gorm.DB {
	var db *gorm.DB
	if db == nil {
		db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		if err != nil {
			fmt.Println(err)
		}
		return db
	}
	return db
}
