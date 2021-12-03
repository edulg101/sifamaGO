package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sifamaGO/src/config"
	"sifamaGO/src/dbService"
	"sifamaGO/src/model"
	"sifamaGO/src/selenium"
	"sifamaGO/src/service"
	"sifamaGO/src/util"
	"text/template"
	"time"
)

func HomeGet(w http.ResponseWriter, r *http.Request) {
	var cookie *http.Cookie
	var err error
	var cookieValue string

	cookie, _ = r.Cookie("sifamaGuid")

	if cookie == nil {
		cookieValue := createRandamHash()

		cookie := http.Cookie{
			Name:    "sifamaGuid",
			Value:   cookieValue,
			Expires: time.Now().Add(365 * 24 * time.Hour),
		}

		http.SetCookie(w, &cookie)
		// session = dbService.CreateNewSession(cookieValue)

	} else {
		cookieValue = cookie.Value
	}

	if r.Method == "GET" {

		err = config.GetEnv()
		if err != nil {
			panic(err)
		}

		service.CleanUpDB(cookieValue)
		f, _ := os.Open(util.ROOTPATH)
		files, _ := f.ReadDir(-1)

		var filesArray []model.Folder
		for _, file := range files {
			if file.IsDir() {
				filesArray = append(filesArray, model.Folder{FolderName: file.Name()})
			}
		}

		data := model.HomeModel{
			Folders: filesArray,
		}

		tmpl := template.Must(template.ParseFiles("view/index.html"))

		tmpl.Execute(w, data)
	}

}

func HomePost(w http.ResponseWriter, r *http.Request) {

	var cookie *http.Cookie
	var err error
	var cookieValue string
	var session *model.Session

	cookie, _ = r.Cookie("sifamaGuid")

	if cookie == nil {
		cookieValue := createRandamHash()
		cookie := http.Cookie{
			Name:    "sifamaGuid",
			Value:   cookieValue,
			Expires: time.Now().Add(365 * 24 * time.Hour),
		}
		http.SetCookie(w, &cookie)
		session = dbService.CreateNewSession(cookieValue)
		fmt.Println(cookieValue)

	} else {
		cookieValue = cookie.Value
		session, err = dbService.FindSessionByHash(cookieValue)
		if err != nil {
			session = dbService.CreateNewSession(cookieValue)
		}
	}

	service.CleanUpDB(cookieValue)

	fmt.Println("metodo post entrou")

	var request Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}

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

	err = selenium.ImportSpreadSheet(util.SPREADSHEETPATH, session)

	fmt.Printf("tempo: %v\n", time.Since(start))

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(fmt.Sprintln(err)))
	}

}

func Compact(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	var request Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}

	if request.Compact {
		util.ORIGINIMAGEPATH = filepath.Join(util.ROOTPATH, request.Folder)
		fmt.Println(util.ORIGINIMAGEPATH)
		message, err := service.ResizeAllImagesInFolder(util.ORIGINIMAGEPATH, util.MAXIMAGEWIDTH)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			message = fmt.Sprint(err)
			w.Write([]byte(message))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(message))
			return
		}

	} else {
		w.Write([]byte("ocorreu um erro"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func createRandamHash() string {
	chave := make([]byte, 5)
	if _, err := rand.Read(chave); err != nil {
		panic(err)
	}
	stringBase64 := base64.StdEncoding.EncodeToString(chave)
	return stringBase64
}
