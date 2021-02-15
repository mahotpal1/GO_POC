package main

import (
	"bufio"
	"context"
	"fmt"
	"global"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {

	log.Println("Starting the application")
	serveWeb()
}

var themeName = getThemeName()
var staticPages = populateStaticPages()

func serveWeb() {

	var gorillaRoute = mux.NewRouter()
	global.DBCollection("user").InsertOne(context.Background(), bson.M{"name": "test"})

	gorillaRoute.HandleFunc("/", serveContent)
	gorillaRoute.HandleFunc("/{pageAlias}", serveContent)
	http.HandleFunc("/process", processor)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/js/", serveResource)
	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":8080", nil)

}

func serveContent(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	pageAlias := urlParams["pageAlias"]
	if pageAlias == "" {
		pageAlias = "Main"
	}

	staticPage := staticPages.Lookup(pageAlias + ".html")
	if staticPage == nil {
		staticPage = staticPages.Lookup("404.html")
		w.WriteHeader(404)
	}
	staticPage.Execute(w, nil)
}

func getThemeName() string {
	return "bs4"
}

func populateStaticPages() *template.Template {
	result := template.New("templates")
	templatePaths := new([]string)

	basePath := "../bin/pages"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ := templateFolder.Readdir(-1)

	for _, pathInfo := range templatePathsRaw {
		log.Println(pathInfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
	}
	result.ParseFiles(*templatePaths...)
	return result
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "../bin/public/" + themeName + req.URL.Path
	var contentType string

	if strings.HasSuffix(path, ".png") {
		contentType = "image/png; charset=utf-8"
	} else if strings.HasSuffix(path, ".jpg") {
		contentType = "image/jpg; charset=utf-8"
	} else if strings.HasSuffix(path, ".css") {
		contentType = "text/css; charset=utf-8"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript; charset=utf-8"
	} else {
		contentType = "text/plain; charset=utf-8"
	}

	log.Println(path)

	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		w.Header().Add("Content-Type", contentType)
		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}

func processor(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	empid := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")

	d := struct {
		Empid    string
		Password string
		Role     string
	}{
		Empid:    empid,
		Password: password,
		Role:     role,
	}
	t, err := template.ParseFiles("../bin/pages/processor.html")

	if err != nil {
		fmt.Println(err)
	}

	t.Execute(w, d)
}

func processor2(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	empid := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")

	d := struct {
		Empid    string
		Password string
		Role     string
	}{
		Empid:    empid,
		Password: password,
		Role:     role,
	}
	t, err := template.ParseFiles("../bin/pages/processor.html")

	if err != nil {
		fmt.Println(err)
	}

	t.Execute(w, d)
}
