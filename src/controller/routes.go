package controller

import (
	"net/http"
	"path/filepath"
	"sifamaGO/src/util"

	"github.com/gorilla/mux"
)

func LoadControllers(r *mux.Router) {
	r.HandleFunc("/", HomeGet).Methods("GET")
	r.HandleFunc("/", HomePost).Methods("POST")

	r.HandleFunc("/report", Report)

	r.HandleFunc("/compact", Compact).Methods("POST")

	r.HandleFunc("/favicon.ico", FaviconHandler)

	r.HandleFunc("/map", Map).GetMethods()

	r.HandleFunc("/checkSifama", checkSifama).Methods("POST")

	staticDir := util.OUTPUTIMAGEFOLDER
	staticURL := "/fotos/"
	r.PathPrefix(staticURL).Handler(http.StripPrefix(staticURL, http.FileServer(http.Dir(staticDir))))

	// Logo Router and static images
	imgesStaticDir := filepath.Join("view", "images")
	cssStaticPath := filepath.Join("view", "css")
	scriptStaticPath := filepath.Join("view", "script")

	// a principio somente para o logo
	imagesStaticURL := "/images/"
	r.PathPrefix(imagesStaticURL).Handler(http.StripPrefix(imagesStaticURL, http.FileServer(http.Dir(imgesStaticDir))))
	// Css Router

	cssUrlPath := "/css/"
	r.PathPrefix(cssUrlPath).Handler(http.StripPrefix(cssUrlPath, http.FileServer(http.Dir(cssStaticPath))))
	// Script Router

	scriptUrlPath := "/script/"
	r.PathPrefix(scriptUrlPath).Handler(http.StripPrefix(scriptUrlPath, http.FileServer(http.Dir(scriptStaticPath))))

}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "view/images/favicon.ico")
}
