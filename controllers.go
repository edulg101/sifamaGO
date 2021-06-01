package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
}

func Report(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		fmt.Println("metodo post entrou")
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.

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
			InicioDigitacao()
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

		tmpl := template.Must(template.ParseFiles("src/template/index.html"))

		tmpl.Execute(w, data)
	}
	if r.Method == "POST" {
		fmt.Println("metodo post entrou")

		folder := r.FormValue("folderselect")

		util.ORIGINIMAGEPATH = filepath.Join(util.ROOTPATH, folder)

		title := r.FormValue("titulo")
		if title != "" {
			util.TITLE = title
		} else {
			today := time.Now().Format("02/01/2006")
			util.TITLE = "Tros Emitidos em " + today
		}

		populateFotosOnDB(util.ORIGINIMAGEPATH)
		err := ImportSpreadSheet(util.SPREADSHEETPATH)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(fmt.Sprintln(err)))
		}

		http.Redirect(w, r, "/report", http.StatusSeeOther)

	}

}

func reportGet(w http.ResponseWriter) {

	var tro Tro
	tros := tro.FindAll()

	totalTro := len(tros)
	totalFotos := 0
	for _, tro := range tros {
		totalFotos += len(tro.Locais)
	}

	data := TroModel{
		Title:      util.TITLE,
		Tro:        tros,
		TotalTro:   totalTro,
		TotalFotos: totalFotos,
	}

	tmpl := template.Must(template.ParseFiles("src/template/report.html"))

	tmpl.Execute(w, data)
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "src/template/images/favicon.ico")
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

	tmpl := template.Must(template.ParseFiles("src/template/map.html"))

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

// func Home1(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		db.CleanUpDB(db.GetDB())
// 		f, _ := os.Open(util.ROOTPATH)
// 		files, _ := f.ReadDir(-1)

// 		var filesArray []Folder
// 		for _, file := range files {
// 			if file.IsDir() {
// 				filesArray = append(filesArray, Folder{FolderName: file.Name()})
// 			}
// 		}

// 		data := HomeModel{
// 			Folders: filesArray,
// 		}

// 		var assetData embed.FS

// 		tmpl, err := template.ParseFS(assetData, "template/index.html")
// 		fmt.Println(err)

// 		tmpl.Execute(w, data)
// 	}
// 	if r.Method == "POST" {
// 		fmt.Println("metodo post entrou")

// 		folder := r.FormValue("folderselect")

// 		util.ORIGINIMAGEPATH = filepath.Join(util.ROOTPATH, folder)

// 		title := r.FormValue("titulo")
// 		if title != "" {
// 			util.TITLE = title
// 		} else {
// 			today := time.Now().Format("02/01/2006")
// 			util.TITLE = "Tros Emitidos em " + today
// 		}

// 		populateFotosOnDB(util.ORIGINIMAGEPATH)
// 		err := ImportSpreadSheet(util.SPREADSHEETPATH)
// 		if err != nil {
// 			w.WriteHeader(http.StatusNotAcceptable)
// 			w.Write([]byte(fmt.Sprintln(err)))
// 		}

// 		http.Redirect(w, r, "/report", http.StatusSeeOther)

// 	}

// }
