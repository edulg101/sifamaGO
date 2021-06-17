package db

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func ConectDB() {
	var err error
	databaseName := ":memory:"
	// databaseName := "gbguru.db"

	err = os.Remove(databaseName)
	if err != nil {
		fmt.Println(err)
	}

	db, err = gorm.Open(sqlite.Open(databaseName), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

}
