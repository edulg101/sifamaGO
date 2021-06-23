package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sifamaGO/src/dbService"
	"sifamaGO/src/model"
	"sifamaGO/src/selenium"
	"sifamaGO/src/service"
	"sifamaGO/src/util"
	"text/template"
)

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
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprint(err)))
		}
		if request.StartDigitacao {
			user := request.User
			password := request.Passd
			returnMessagem, err := selenium.InicioDigitacao(user, password)
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				w.Write([]byte(fmt.Sprint(err)))
				return
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(returnMessagem))
				return
			}
		}
		if request.Restart {
			err := restart(w, r)
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				w.Write([]byte(fmt.Sprintln(err)))
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	}
	if r.Method == "GET" {
		reportGet(w, r)
	}

}
func restart(w http.ResponseWriter, r *http.Request) error {

	cookie, _ := r.Cookie("sifamaGuid")
	cookieValue := cookie.Value

	session, err := dbService.FindSessionByHash(cookieValue)
	if err != nil {
		fmt.Println(err)
	}

	service.CleanUpDB(cookieValue)
	fmt.Println("output image folder:", util.OUTPUTIMAGEFOLDER)

	// service.PopulateFotosOnDB(util.ORIGINIMAGEPATH, cookie.Value)
	err = selenium.ImportSpreadSheet(util.SPREADSHEETPATH, session)
	if err != nil {
		return err
	}

	return nil

}

func reportGet(w http.ResponseWriter, r *http.Request) {

	cookie, _ := r.Cookie("sifamaGuid")
	cookieValue := cookie.Value

	tros, err := service.FindAllBySession(cookieValue)

	if err != nil {
		panic(err)
	}

	var localWithNoFotos []model.Local

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

	data := model.TroModel{
		Title:            util.TITLE,
		Tro:              tros,
		TotalTro:         totalTro,
		TotalFotos:       totalFotos,
		LocalWithNoFotos: localWithNoFotos,
	}

	tmpl := template.Must(template.ParseFiles("view/report.html"))

	tmpl.Execute(w, data)
}
