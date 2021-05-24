package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"gorm.io/gorm"
)

func errorHandle(e error) {
	if e != nil {
		log.Println(e)
		panic(e)
	}
}

func properTitle(input string) string {
	words := strings.Fields(input)

	for index, word := range words {
		if len(word) < 4 {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}

func checkForNameSize(fullPath string) string {

	//check if FileName is too large
	fileName := filepath.Base(fullPath)
	path := filepath.Dir(fullPath)

	nameLength := len(fileName)
	if nameLength > 90 {
		oldNameFile := fileName

		fmt.Println("reduzindo o tamanho do nome do arquivo")
		// bytes := []byte(fileName)
		half := nameLength / 2
		extraHalf := (nameLength - 90) / 2
		// bytes = append(bytes[:half-extraHalf], bytes[(half+extraHalf+1):]...)
		newFileName := oldNameFile[:half-extraHalf] + oldNameFile[half+extraHalf:]

		// newDecoder não necessário para Windows
		// fileName, err := charmap.CodePage850.NewDecoder().String(fileName)
		oldPath := filepath.Join(path, oldNameFile)
		newPath := filepath.Join(path, newFileName)
		fileName = newFileName
		err := os.Rename(oldPath, newPath)
		errorHandle(err)

	}
	return fileName
}

func keepMouseMoving() {
	for {
		robotgo.MoveMouse(100, 300)
		time.Sleep(time.Minute * 2)
		robotgo.MoveMouse(300, 500)
		time.Sleep(time.Minute * 2)
	}
}

func isTrechosDNIT(km float64) bool {

	return (km >= 221.29 && km <= 230.06) || (km >= 277 && km <= 360)
}

func isLocationValid(caption string, local *Local) string {

	kmInicial := local.KmInicialDouble
	kmFinal := local.KmFinalDouble
	palavraChave := local.Tro.PalavraChave
	oldKmInicial := local.KmInicial
	oldKmFinal := local.KmFinal

	if strings.Contains(local.Rodovia, "364") && strings.Contains(strings.ToLower(local.Sentido), "decrescente") {
		if (kmInicial > 0 && kmInicial < 20) || (kmFinal > 0 && kmFinal < 20) {
			newKmInicial, newKmFinal := interpolationLocal(local)
			kmFinalStr := fmt.Sprintf("%.3f", newKmFinal)
			kmInicialStr := fmt.Sprintf("%.3f", newKmInicial)
			local.KmFinal = kmFinalStr
			local.KmFinalDouble = newKmFinal
			local.KmInicial = kmInicialStr
			local.KmInicialDouble = newKmInicial
			caption = caption + " (km da 364 Variante : " + oldKmInicial + " - " + oldKmFinal + " )"
			db.Save(&local)
		}

	}

	kmInicial = local.KmInicialDouble
	kmFinal = local.KmFinalDouble

	checkkmInicial, dnit := checkKm(local.Rodovia, kmInicial, palavraChave)
	checkkmFinal, dnit1 := checkKm(local.Rodovia, kmFinal, palavraChave)

	local.Valid = checkkmInicial && checkkmFinal
	local.TrechoDNIT = dnit || dnit1

	db.Save(&local)

	return caption
}

func interpolationLocal(local *Local) (float64, float64) {

	kmInicial := local.KmInicialDouble
	kmFinal := local.KmFinalDouble
	newKmInicial := interpolationKm(kmInicial)
	newKmFinal := interpolationKm(kmFinal)
	return newKmInicial, newKmFinal
}
func interpolationKm(km float64) float64 {
	if km <= 10.94 {
		return 360 - ((km - 1.0) / 9.94 * 7.7)
	} else if km <= 12.25 {
		return 351 + (12.25 - km)
	} else if km <= 19.67 {
		return 351 - ((km-12.25)/7.42)*7.2
	} else {
		return km
	}

}

func checkKm(rodovia string, km float64, palavraChave string) (bool, bool) {

	disposicao := getDisposicaoLegal(palavraChave)
	edificacoes := false

	if disposicao[0] == "5" && disposicao[1] == "742" {
		edificacoes = true
	}

	switch rodovia {
	case "70":
		if km < 495.9 || km > 524 {
			return false, false
		}
		return true, false

	case "163":
		return (km < 120 || (km <= 855 && km >= 507.1)), false

	case "364":
		if km < 201 || km > 588.2 {
			return false, false
		} else if km <= 211.3 {
			return true, false
		} else if !edificacoes && isTrechosDNIT(km) {
			return true, true
		} else if edificacoes && isTrechosDNIT(km) {
			return true, false
		}

	default:
		return false, false
	}
	return false, false
}
func cleanUpDB(db *gorm.DB) {
	rows := db.Exec("DELETE FROM 'fotos' WHERE id > 0")
	fmt.Println("fotos deletadas: ", rows.RowsAffected)
	rows = db.Exec("DELETE FROM 'locals' WHERE id > 0")
	fmt.Println("locals deletadas: ", rows.RowsAffected)
	rows = db.Exec("DELETE FROM 'tros' WHERE id > 0")
	fmt.Println("tros deletadas: ", rows.RowsAffected)
}
