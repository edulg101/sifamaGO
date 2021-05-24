package main

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Foto struct {
	gorm.Model
	ID      uint
	Nome    string
	Path    template.URL
	Legenda string
	LocalID uint
	Local   Local
	UrlPath template.URL
}

type Local struct {
	gorm.Model
	ID               uint
	NumIdentificacao string
	Data             string
	Hora             string
	Rodovia          string
	Pista            string
	KmInicial        string
	KmFinal          string
	Sentido          string
	KmInicialDouble  float64
	KmFinalDouble    float64
	TrechoDNIT       bool
	Valid            bool
	Fotos            []Foto `gorm:"ForeignKey:LocalID"`
	TroID            uint
	Tro              Tro
}

type Tro struct {
	gorm.Model
	ID            uint
	PalavraChave  string
	Observacao    string
	Prazo         string
	TipoPrazo     string
	Severidade    string
	Disposicao    string
	DisposicaoCod string
	DisposicaoArt string
	Locais        []Local // `gorm:"ForeignKey:TroID"`
}

type Folder struct {
	FolderName string
}

func (tro Tro) findAll() []Tro {
	var troList []Tro
	db.Preload("Locais.Fotos").Preload(clause.Associations).Find(&troList)
	return troList
}
func (local Local) findAllLocais() []Local {
	var localList []Local
	db.Preload("Fotos").Find(&localList)
	return localList
}
func (tro Tro) findAllFotos() []Tro {
	var troList []Tro
	db.Preload("Locais.Fotos").Preload(clause.Associations).Find(&troList)
	return troList
}

type FilesModel struct {
	Folders []Folder
}

type TroModel struct {
	Data     string
	Tro      []Tro
	TotalTro int
	Folders  []Folder
}

// var indexHTML embed.FS

func main() {

	conectDB()
	// populateFotosOnDB(ORIGINIMAGEPATH)
	// _, err := importSpreadSheet(SPREADSHEETPATH)
	// errorHandle(err)

	r := mux.NewRouter()

	r.HandleFunc("/", home)

	r.HandleFunc("/inicial", inicial)

	r.HandleFunc("/favicon.ico", faviconHandler)

	// Image server
	// staticDir := "template\\images\\"
	staticDir := OUTPUTIMAGEFOLDER
	staticURL := "/fotos/"
	r.PathPrefix(staticURL).Handler(http.StripPrefix(staticURL, http.FileServer(http.Dir(staticDir))))

	// Logo Router and static images

	imgesStaticDir := "template\\images\\"
	imagesStaticURL := "/images/"
	r.PathPrefix(imagesStaticURL).Handler(http.StripPrefix(imagesStaticURL, http.FileServer(http.Dir(imgesStaticDir))))

	// Css Router
	cssStaticPath := "template\\css\\"
	cssUrlPath := "/css/"
	r.PathPrefix(cssUrlPath).Handler(http.StripPrefix(cssUrlPath, http.FileServer(http.Dir(cssStaticPath))))

	// Script Router

	scrpitStaticPath := "template\\script\\"
	scriptUrlPath := "/script/"
	r.PathPrefix(scriptUrlPath).Handler(http.StripPrefix(scriptUrlPath, http.FileServer(http.Dir(scrpitStaticPath))))

	err := http.ListenAndServe(":8080", r)
	errorHandle(err)

	time.Sleep(time.Minute)

}
func restart() {

	cleanUpDB(db)

	// populateFotosOnDB(ORIGINIMAGEPATH)
	// _, err := importSpreadSheet(SPREADSHEETPATH)
	// errorHandle(err)

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
