package main

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"

	"sifamaGO/db"
	"sifamaGO/util"
)

// var indexHTML embed.FS

func main() {

	db.ConectDB()

	db.GetDB().AutoMigrate(&Tro{})
	db.GetDB().AutoMigrate(&Local{})
	db.GetDB().AutoMigrate(&Foto{})
	// populateFotosOnDB(ORIGINIMAGEPATH)
	// _, err := importSpreadSheet(SPREADSHEETPATH)
	// errorHandle(err)

	r := mux.NewRouter()

	r.HandleFunc("/", Home)

	r.HandleFunc("/report", Report)

	r.HandleFunc("/favicon.ico", FaviconHandler)

	staticDir := util.OUTPUTIMAGEFOLDER
	staticURL := "/fotos/"
	r.PathPrefix(staticURL).Handler(http.StripPrefix(staticURL, http.FileServer(http.Dir(staticDir))))

	// Logo Router and static images

	imgesStaticDir := filepath.Join("template", "images")
	cssStaticPath := filepath.Join("template", "css")
	scriptStaticPath := filepath.Join("template", "script")

	imagesStaticURL := "/images/"
	r.PathPrefix(imagesStaticURL).Handler(http.StripPrefix(imagesStaticURL, http.FileServer(http.Dir(imgesStaticDir))))
	// Css Router
	cssUrlPath := "/css/"
	r.PathPrefix(cssUrlPath).Handler(http.StripPrefix(cssUrlPath, http.FileServer(http.Dir(cssStaticPath))))
	// Script Router
	scriptUrlPath := "/script/"
	r.PathPrefix(scriptUrlPath).Handler(http.StripPrefix(scriptUrlPath, http.FileServer(http.Dir(scriptStaticPath))))

	err := http.ListenAndServe(":8080", r)
	panic(err)

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
