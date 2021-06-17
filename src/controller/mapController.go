package controller

import (
	"fmt"
	"net/http"
	"sifamaGO/src/model"
	"sifamaGO/src/tests/geo"
	"text/template"
)

func Map(w http.ResponseWriter, r *http.Request) {

	var pontos []model.Geolocation
	var locations []model.Geolocation
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
	Points []model.Geolocation
}

func errorHandle(err error) {
	fmt.Println(err)
}
