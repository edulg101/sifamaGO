package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/zserge/lorca"

	"sifamaGO/src/config"
	"sifamaGO/src/db"
	"sifamaGO/src/util"
)

var PORT string

func main() {

	// util.ToAscII("D:\\sifamadocs\\inPics\\input")

	ui, er := lorca.New("", "", 1000, 800)
	if er != nil {
		log.Fatal(er)
	}
	defer ui.Close()

	db.ConectDB()

	db.GetDB().AutoMigrate(&Tro{})
	db.GetDB().AutoMigrate(&Local{})
	db.GetDB().AutoMigrate(&Foto{})

	r := mux.NewRouter()

	r.HandleFunc("/", Home)

	r.HandleFunc("/report", Report)

	r.HandleFunc("/favicon.ico", FaviconHandler)

	r.HandleFunc("/map", Map).GetMethods()

	staticDir := util.OUTPUTIMAGEFOLDER
	staticURL := "/fotos/"
	r.PathPrefix(staticURL).Handler(http.StripPrefix(staticURL, http.FileServer(http.Dir(staticDir))))

	// Logo Router and static images
	imgesStaticDir := filepath.Join("src", "template", "images")
	cssStaticPath := filepath.Join("src", "template", "css")
	scriptStaticPath := filepath.Join("src", "template", "script")

	// a principio somente para o logo
	imagesStaticURL := "/images/"
	r.PathPrefix(imagesStaticURL).Handler(http.StripPrefix(imagesStaticURL, http.FileServer(http.Dir(imgesStaticDir))))
	// Css Router

	cssUrlPath := "/css/"
	r.PathPrefix(cssUrlPath).Handler(http.StripPrefix(cssUrlPath, http.FileServer(http.Dir(cssStaticPath))))
	// Script Router

	scriptUrlPath := "/script/"
	r.PathPrefix(scriptUrlPath).Handler(http.StripPrefix(scriptUrlPath, http.FileServer(http.Dir(scriptStaticPath))))

	PORT, _, _ = config.GetEnv()

	go callUI(ui)

	fmt.Println("porta", PORT)

	err := http.ListenAndServe(":"+PORT, r)
	if err != nil {
		panic(err)
	}

}

func callUI(ui lorca.UI) {
	time.Sleep(time.Second / 2)
	add := "http://localhost:" + PORT
	fmt.Println(PORT)
	ui.Load(add)
}

// tmpl := template.Must(template.ParseFiles("template/index.html"))

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == "POST" {
// 			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
// 			if err := r.ParseForm(); err != nil {
// 				fmt.Fprintf(w, "ParseForm() err: %v", err)
// 				return
// 			}
// 			todo := r.FormValue("todo")
// 			db.Create(&Todo{Title: todo, Done: false})
// 		}
// 		//Request not POST
// 		var todos []Todo
// 		db.Find(&todos)
// 		data := TodoPageData{
// 			PageTitle: "Lista de Tarefas",
// 			Todos:     todos,
// 		}
// 		tmpl.Execute(w, data)

// 	})

// 	http.HandleFunc("/done/", func(w http.ResponseWriter, r *http.Request) {
// 		id := strings.TrimPrefix(r.URL.Path, "/done/")
// 		var todo Todo
// 		db.First(&todo, id)
// 		todo.Done = true
// 		db.Save(&todo)
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 	})

// 	http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
// 		id := strings.TrimPrefix(r.URL.Path, "/delete/")
// 		db.Delete(&Todo{}, id)
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 	})
