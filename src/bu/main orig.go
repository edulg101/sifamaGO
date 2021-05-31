package backup

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/zserge/lorca"

	"sifamaGO/db"
	"sifamaGO/util"
)

//go:embed template/*
var assetData embed.FS

func main1() {

	ui, err := lorca.New("", "", 1000, 800)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

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

	text, e := assetData.Open("template/css/style.css")
	fmt.Println(e)
	info, e := text.Stat()
	fmt.Println(e)
	fmt.Println(info.Name())

	imagesStaticURL := "/images/"
	r.PathPrefix(imagesStaticURL).Handler(http.StripPrefix(imagesStaticURL, http.FileServer(http.Dir(imgesStaticDir))))
	// Css Router

	cssUrlPath := "/css/"
	r.PathPrefix(cssUrlPath).Handler(http.StripPrefix(cssUrlPath, http.FileServer(http.Dir(cssStaticPath))))
	// Script Router

	scriptUrlPath := "/script/"
	r.PathPrefix(scriptUrlPath).Handler(http.StripPrefix(scriptUrlPath, http.FileServer(http.Dir(scriptStaticPath))))

	go http.ListenAndServe(":8080", r)
	fmt.Println("entrou")

	time.Sleep(time.Second)
	ui.Load("http://localhost:8080/")

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")
	// panic(err)

}
func Home1(w http.ResponseWriter, r *http.Request) {
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

		tmpl, err := template.ParseFS(assetData, "template/index.html")
		fmt.Println(err)

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

tmpl := template.Must(template.ParseFiles("template/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			todo := r.FormValue("todo")
			db.Create(&Todo{Title: todo, Done: false})
		}
		//Request not POST
		var todos []Todo
		db.Find(&todos)
		data := TodoPageData{
			PageTitle: "Lista de Tarefas",
			Todos:     todos,
		}
		tmpl.Execute(w, data)

	})

	http.HandleFunc("/done/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/done/")
		var todo Todo
		db.First(&todo, id)
		todo.Done = true
		db.Save(&todo)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/delete/")
		db.Delete(&Todo{}, id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
