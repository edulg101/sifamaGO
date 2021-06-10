package config

import (
	"fmt"
	"log"
	"os"
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
	util.PWD = os.Getenv("PASSW")
	util.USER = os.Getenv("USS")
	util.SELENIUMPATH = os.Getenv("SELENIUMPATH")
	util.SPREADSHEETPATH = os.Getenv("SPREADSHEETPATH")
	util.FONTPATH = os.Getenv("FONTPATH")
	util.OUTPUTIMAGEFOLDER = os.Getenv("OUTPUTIMAGEFOLDER")

	MaxImageWidth, err := strconv.ParseUint(os.Getenv("MAX-IMAGE-WIDTH"), 10, 32)
	if err != nil {
		return fmt.Errorf("não foi possível identificar a constante: MAX-IMAGE-WIDTH")
	}
	util.MAXIMAGEWIDTH = uint(MaxImageWidth)
	util.ROOTPATH = os.Getenv("ROOTPATH")

	return nil
}
