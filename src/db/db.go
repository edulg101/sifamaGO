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
	// databaseName := "gb9.db"

	err = os.Remove(databaseName)
	if err != nil {
		fmt.Println(err)
	}

	db, err = gorm.Open(sqlite.Open(databaseName), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

}

func CleanUpDB(db *gorm.DB) {
	rows := db.Exec("DELETE FROM 'fotos' WHERE id > 0")
	fmt.Println("fotos deletadas: ", rows.RowsAffected)
	rows = db.Exec("DELETE FROM 'locals' WHERE id > 0")
	fmt.Println("locals deletadas: ", rows.RowsAffected)
	rows = db.Exec("DELETE FROM 'tros' WHERE id > 0")
	fmt.Println("tros deletadas: ", rows.RowsAffected)
}
