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

	util.PWD = os.Getenv("PWD")

	util.USER = os.Getenv("USER")

	MaxImageWidth, err := strconv.ParseUint(os.Getenv("MAX-IMAGE-WIDTH"), 10, 32)
	if err != nil {
		return fmt.Errorf("não foi possível identificar a constante: MAX-IMAGE-WIDTH")
	}
	util.MAXIMAGEWIDTH = uint(MaxImageWidth)
	util.ROOTPATH = os.Getenv("ROOTPATH")

	return nil
}
