package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sifamaGO/src/util"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	StringConexaoBanco = ""
	Porta              = 0
)

func GetEnv() error {
	var err error
	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	util.PORT = os.Getenv("PORT")
	util.SELENIUMPATH = os.Getenv("SELENIUMPATH")
	util.SPREADSHEETPATH = os.Getenv("SPREADSHEETPATH")
	util.FONTPATH = os.Getenv("FONTPATH")
	util.OUTPUTIMAGEFOLDER = os.Getenv("OUTPUTIMAGEFOLDER")
	util.EXIFTOOL = os.Getenv("EXIFTOOL")
	util.CHECKSPREADSHEETPATH = os.Getenv("CHECKSPREADSHEETPATH")

	MaxImageWidth, err := strconv.ParseUint(os.Getenv("MAX-IMAGE-WIDTH"), 10, 32)
	if err != nil {
		return fmt.Errorf("não foi possível identificar a constante: MAX-IMAGE-WIDTH")
	}
	util.MAXIMAGEWIDTH = uint(MaxImageWidth)
	util.ROOTPATH = os.Getenv("ROOTPATH")

	// verify if imagefolder is relative or absolute path:

	currentDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(currentDirectory, util.OUTPUTIMAGEFOLDER)

	let1 := string(util.OUTPUTIMAGEFOLDER[1])
	let0 := string(util.OUTPUTIMAGEFOLDER[0])

	if !(let1 == ":" || let0 == "/") {
		util.OUTPUTIMAGEFOLDER = fullPath
	}
	fmt.Println("util.OUTPUTIMAGEFOLDER --> ", util.OUTPUTIMAGEFOLDER)
	return nil
}
