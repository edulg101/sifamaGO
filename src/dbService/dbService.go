package dbService

import (
	"fmt"
	"sifamaGO/src/db"
	"sifamaGO/src/model"

	"gorm.io/gorm"
)

func FindSessionByHash(hash string) (*model.Session, error) {

	var sessions []model.Session
	db.GetDB().Find(&sessions)
	for _, session := range sessions {
		if session.Hash == hash {
			return &session, nil
		}
	}
	return nil, fmt.Errorf("nao encontrou.")
}

func CreateNewSession(hash string) *model.Session {
	var session model.Session
	session.Hash = hash
	db.GetDB().Create(&session)

	return &session
}

func GetDB() *gorm.DB {
	return db.GetDB()
}
