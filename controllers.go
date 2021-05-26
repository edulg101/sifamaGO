package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sifamaGO/db"
	"sifamaGO/util"
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
			restart()
			w.WriteHeader(http.StatusOK)
		}
	}
	if r.Method == "GET" {
		reportGet(w)
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

	tmpl := template.Must(template.ParseFiles("template/report.html"))

	tmpl.Execute(w, data)
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

		tmpl := template.Must(template.ParseFiles("template/index.html"))

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
		_, err := ImportSpreadSheet(util.SPREADSHEETPATH)
		errorHandle(err)
		http.Redirect(w, r, "/report", http.StatusSeeOther)

	}

}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/images/favicon.ico")
}

func restart() {

	db.CleanUpDB(db.GetDB())
	fmt.Println("output image folder:", util.OUTPUTIMAGEFOLDER)

	populateFotosOnDB(util.ORIGINIMAGEPATH)
	_, err := ImportSpreadSheet(util.SPREADSHEETPATH)
	errorHandle(err)

	// populateFotosOnDB(ORIGINIMAGEPATH)
	// _, err := importSpreadSheet(SPREADSHEETPATH)
	// errorHandle(err)

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
