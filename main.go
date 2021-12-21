package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/zserge/lorca"

	"sifamaGO/src/config"
	"sifamaGO/src/controller"
	"sifamaGO/src/db"
	"sifamaGO/src/model"
	"sifamaGO/src/util"
)

func testPaths() {

	currentDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	imgPath1 := filepath.Join(currentDirectory, util.OUTPUTIMAGEFOLDER)

	imgpath := util.OUTPUTIMAGEFOLDER

	_, err = os.Stat(imgpath)
	if os.IsNotExist(err) {
		imgpath = imgPath1
	}

	fmt.Println(imgPath1)
	fmt.Println(imgpath)

}

func removeDots(filename string) string {
	count := strings.Count(filename, ".")
	if count > 1 {
		filename = strings.Replace(filename, ".", "", count-1)
	}
	return filename
}

func startRemoveDots() {
	path := "D:\\OneDrive - ANTT- Agencia Nacional de Transportes Terrestres\\SharePointCRO\\RTA\\2 - Diários\\2021-11\\2021_11_17 RDO"
	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)

		}
		fmt.Println("currentpath: ", currentPath)
		dir, name := filepath.Split(currentPath)
		name = removeDots(name)

		err = os.Rename(currentPath, filepath.Join(dir, name))
		fmt.Printf("renomeando arquivo: %s", currentPath)

		return err
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("entrou")
}

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
