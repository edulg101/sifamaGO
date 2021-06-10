package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sifamaGO/src/config"
	"sifamaGO/src/db"
	"sifamaGO/src/util"
	"sifamaGO/src/util/geo"
	"text/template"
	"time"

	"gorm.io/gorm/clause"
)

type Request struct {
	StartDigitacao bool
	Restart        bool
	Compact        bool
	Folder         string
	Title          string
	User           string
	Passd          string
}

func Report(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		fmt.Println("metodo post entrou")

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		var request Request
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		errorHandle(err)

		fmt.Println("start digitacao", request.StartDigitacao)
		fmt.Println("Restart", request.Restart)
		if request.StartDigitacao {
			util.USER = request.User
			util.PWD = request.Passd
			err := InicioDigitacao()
			if err != nil {

				w.WriteHeader(http.StatusNotAcceptable)
				w.Write([]byte(fmt.Sprint(err)))
				return
			}
		}
		if request.Restart {
			err := restart()
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				w.Write([]byte(fmt.Sprintln(err)))
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	}
	if r.Method == "GET" {
		reportGet(w)
	}

}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err1 := config.GetEnv()
		if err1 != nil {
			panic(err1)
		}
		db.CleanUpDB(db.GetDB())
		f, _ := os.Open(util.ROOTPATH)
		files, _ := f.ReadDir(-1)

		var filesArray []Folder
		for _, file := range files {
			if file.IsDir() {
				filesArray = append(filesArray, Folder{FolderName: file.Name()})
			}
		}

		data := HomeModel{
			Folders: filesArray,
		}

		tmpl := template.Must(template.ParseFiles("view/index.html"))

		tmpl.Execute(w, data)
	}
	if r.Method == "POST" {
		fmt.Println("metodo post entrou")

		var request Request
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		errorHandle(err)

		fmt.Println("start digitacao", request.StartDigitacao)
		fmt.Println("Restart", request.Restart)

		folder := request.Folder

		util.ORIGINIMAGEPATH = filepath.Join(util.ROOTPATH, folder)

		title := request.Title

		if title != "" {
			util.TITLE = title
		} else {
			today := time.Now().Format("02/01/2006")
			util.TITLE = "Tros Emitidos em " + today
		}

		start := time.Now()

		err = ImportSpreadSheet(util.SPREADSHEETPATH)

		fmt.Printf("tempo: %v\n", time.Since(start))

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(fmt.Sprintln(err)))
		}

		http.Redirect(w, r, "/report", http.StatusSeeOther)

	}

}

func reportGet(w http.ResponseWriter) {

	var tro Tro
	tros := tro.FindAll()

	var localWithNoFotos []Local

	totalTro := len(tros)
	totalFotos := 0
	for _, tro := range tros {
		totalFotos += len(tro.Locais)
		locais := tro.Locais
		for _, loc := range locais {
			if len(loc.Fotos) < 1 {
				localWithNoFotos = append(localWithNoFotos, loc)
			}
		}

	}

	data := TroModel{
		Title:            util.TITLE,
		Tro:              tros,
		TotalTro:         totalTro,
		TotalFotos:       totalFotos,
		LocalWithNoFotos: localWithNoFotos,
	}

	tmpl := template.Must(template.ParseFiles("view/report.html"))

	tmpl.Execute(w, data)
}

func Compact(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	var request Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	errorHandle(err)

	if request.Compact {
		util.ORIGINIMAGEPATH = filepath.Join(util.ROOTPATH, request.Folder)
		fmt.Println(util.ORIGINIMAGEPATH)
		err := ResizeAllImagesInFolder(util.ORIGINIMAGEPATH, util.MAXIMAGEWIDTH)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			message := fmt.Sprint(err)
			fmt.Println(message)
			w.Write([]byte(message))
		} else {
			w.WriteHeader(http.StatusOK)
		}

	} else {
		w.Write([]byte("ocorreu um erro"))
		w.WriteHeader(http.StatusBadRequest)
	}

}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "view/images/favicon.ico")
}
func Map(w http.ResponseWriter, r *http.Request) {

	var pontos []Geolocation
	var locations []Geolocation
	geo.GetDBGEO().Find(&locations)

	for _, location := range locations {
		km := location.Km
		res := km - float64(int(km))
		if res == 0 {
			pontos = append(pontos, location)
		}
	}

	data := Pontos{
		Points: pontos,
	}

	tmpl := template.Must(template.ParseFiles("view/map.html"))

	tmpl.Execute(w, data)
}

type Pontos struct {
	Points []Geolocation
}

func restart() error {

	db.CleanUpDB(db.GetDB())
	fmt.Println("output image folder:", util.OUTPUTIMAGEFOLDER)

	populateFotosOnDB(util.ORIGINIMAGEPATH)
	err := ImportSpreadSheet(util.SPREADSHEETPATH)
	if err != nil {
		return err
	}

	return nil

}

func (tro Tro) FindAll() []Tro {
	var troList []Tro
	db.GetDB().Preload("Locais.Fotos").Preload(clause.Associations).Find(&troList)

	return troList
}
func (local Local) FindAll() []Local {
	var localList []Local
	db.GetDB().Preload("Fotos").Find(&localList)
	return localList
}
func (tro Tro) findAllFotos() []Tro {
	var troList []Tro
	db.GetDB().Preload("Locais.Fotos").Preload(clause.Associations).Find(&troList)
	return troList
}
