package service

import (
	"sifamaGO/src/db"
	"sifamaGO/src/dbService"
	"sifamaGO/src/model"

	"gorm.io/gorm/clause"
)

func FindAllTro() []model.Tro {
	var troList []model.Tro
	dbService.GetDB().Preload("Locais.Fotos").Preload(clause.Associations).Find(&troList)

	return troList
}

func FindAllBySession(hash string) ([]model.Tro, error) {
	var troList []model.Tro
	var err error
	var session *model.Session
	dbService.GetDB().Preload("Locais.Fotos").Preload(clause.Associations).Find(&troList)
	session, err = dbService.FindSessionByHash(hash)
	if err != nil {
		return nil, err
	}

	var troListResult []model.Tro
	for _, tro := range troList {
		if tro.SessionID == session.ID {
			troListResult = append(troListResult, tro)
		}
	}

	return troListResult, err
}

func FindAllLocal() []model.Local {
	var localList []model.Local
	dbService.GetDB().Preload("Fotos").Find(&localList)
	return localList
}
func findAllFotos() []model.Foto {
	var photoList []model.Foto
	dbService.GetDB().Find(&photoList)
	return photoList
}
func CleanUpDB(hash string) {

	tros, err := FindAllBySession(hash)
	if err != nil {
		return
	}

	// for _, tro := range tros {
	// 	for _, local := range locals {
	// 		if local.TroID == tro.ID {
	// 			db.GetDB().Select(clause.Associations).Delete(&local)
	// 		}
	// 		db.GetDB().Select(clause.Associations).Delete(&tro)
	// 	}
	// }

	var troListToDelete []int
	var localListToDelete []int
	var fotoListToDelete []int

	for _, tro := range tros {
		locals := tro.Locais
		troListToDelete = append(troListToDelete, int(tro.ID))
		for _, local := range locals {
			fotos := local.Fotos
			localListToDelete = append(localListToDelete, int(local.ID))
			for _, foto := range fotos {
				fotoListToDelete = append(fotoListToDelete, int(foto.ID))
			}
		}
	}

	var trotoDelete []model.Tro
	var localToDelete []model.Local
	var fotoToDelete []model.Foto
	db.GetDB().Delete(&fotoToDelete, fotoListToDelete)
	db.GetDB().Delete(&localToDelete, localListToDelete)
	db.GetDB().Delete(&trotoDelete, troListToDelete)

	// db.GetDB().Exec("DELETE FROM 'tros' WHERE session_id = ?", sessionID)

}
