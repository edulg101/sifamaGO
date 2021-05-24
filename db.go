package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func conectDB() {
	var err error
	// databaseName := ":memory:"
	databaseName := "gb9.db"

	err = os.Remove(databaseName)
	if err != nil {
		fmt.Println(err)
	}

	db, err = gorm.Open(sqlite.Open(databaseName), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Tro{})
	db.AutoMigrate(&Local{})
	db.AutoMigrate(&Foto{})

}

type Request struct {
	StartDigitacao bool
	Restart        bool
}

func home(w http.ResponseWriter, r *http.Request) {

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
			inicioDigitacao()
		}
		if request.Restart {
			restart()
			w.WriteHeader(http.StatusOK)

			// http.Redirect(w, r, "/", http.StatusAccepted)
			// homeGet(w, "", "")
		}

	}
	if r.Method == "GET" {
		homeGet(w, "", "")
	}

}

func homeGet(w http.ResponseWriter, titulo, folder string) {
	var tro Tro
	tros := tro.findAll()

	var today string
	totalTro := len(tros)
	if titulo != "" {
		today = titulo
	} else {
		today = time.Now().Format("02/01/2006")
		titulo = "Tros Emitidos em " + today
	}

	data := TroModel{
		Data:     titulo,
		Tro:      tros,
		TotalTro: totalTro,
	}

	tmpl := template.Must(template.ParseFiles("template/index.html"))

	tmpl.Execute(w, data)
}

func inicial(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cleanUpDB(db)
		f, _ := os.Open(ROOTPATH)
		files, _ := f.ReadDir(-1)

		var filesArray []Folder
		for _, file := range files {
			if file.IsDir() {
				filesArray = append(filesArray, Folder{FolderName: file.Name()})
			}
		}

		data := FilesModel{
			Folders: filesArray,
		}

		tmpl := template.Must(template.ParseFiles("template/inicio.html"))

		tmpl.Execute(w, data)
	}
	if r.Method == "POST" {
		fmt.Println("metodo post entrou")

		folder := r.FormValue("folderselect")

		ORIGINIMAGEPATH = filepath.Join(ROOTPATH, folder)

		titulo := r.FormValue("titulo")

		populateFotosOnDB(ORIGINIMAGEPATH)
		_, err := importSpreadSheet(SPREADSHEETPATH)
		errorHandle(err)

		homeGet(w, titulo, folder)
	}

}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/images/favicon.ico")
}
