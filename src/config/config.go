package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	StringConexaoBanco = ""
	Porta              = 0
)

func GetEnv() (PORT, USER, PWD string) {
	var err error
	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	PORT = os.Getenv("PORT")

	PWD = os.Getenv("PWD")

	USER = os.Getenv("USER")
	return
}
