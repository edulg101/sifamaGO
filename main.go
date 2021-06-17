package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zserge/lorca"

	"sifamaGO/src/config"
	"sifamaGO/src/controller"
	"sifamaGO/src/db"
	"sifamaGO/src/model"
	"sifamaGO/src/util"
)

func main() {

	err1 := config.GetEnv()
	if err1 != nil {
		panic(err1)
	}

	ui, er := lorca.New("", "", 1000, 800)
	if er != nil {
		log.Fatal(er)
	}
	defer ui.Close()

	db.ConectDB()

	db.GetDB().AutoMigrate(&model.Session{})
	db.GetDB().AutoMigrate(&model.Tro{})
	db.GetDB().AutoMigrate(&model.Local{})
	db.GetDB().AutoMigrate(&model.Foto{})

	r := mux.NewRouter()
	controller.LoadControllers(r)

	go callUI(ui)

	fmt.Println("porta", util.PORT)

	err := http.ListenAndServe(":"+util.PORT, r)
	if err != nil {
		panic(err)
	}

}

func callUI(ui lorca.UI) {
	time.Sleep(time.Second / 2)
	add := "http://localhost:" + util.PORT
	ui.Load(add)
}
