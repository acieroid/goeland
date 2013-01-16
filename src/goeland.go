package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const templatesDir = "../templates"
const staticDir = "../static"

type Page struct {
	Content template.HTML
}

type Message struct {
	Message   string
	Action    string
	ActionURL string
}

func RenderTemplate(w http.ResponseWriter, tmpl string, arg interface{}) {
	skel, err := template.ParseFiles(fmt.Sprintf("%s/skeleton.html", templatesDir))
	if err != nil {
		log.Println("Error when parsing a template: %s", err)
		return
	}

	content, err := template.ParseFiles(fmt.Sprintf("%s/%s.html", templatesDir, tmpl))
	if err != nil {
		log.Println("Error when parsing a template: %s", err)
		return
	}

	buff := &bytes.Buffer{}
	err = content.Execute(buff, arg)
	if err != nil {
		log.Println("Error when executing a template: %s", err)
		return
	}

	c := template.HTML(buff.String())
	page := &Page{Content: c}
	err = skel.Execute(w, page)
	if err != nil {
		log.Println("Error when executing a template: %s", err)
		return
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index", nil)
}

func About(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "about", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	l := NewList(name)
	l.AddItem(&TodoListItem{1, "Foo", "", "Todo", nil})
	l.GetItem(1).AddItem(&TodoListItem{2, "Bar", "", "Todo", nil})
	l.Save()
	http.Redirect(w, r, "/view/"+l.Id, http.StatusFound)
}

func View(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[6:]
	l := LoadList(id)
	if l != nil {
		RenderTemplate(w, "view", l)
	} else {
		RenderTemplate(w, "message",
			&Message{"This list doesn't exists", "Back", "/"})
	}
}

func Save(w http.ResponseWriter, r *http.Request) {
	l := ParseList([]byte(r.FormValue("list")))
	if l == nil {
		w.Write([]byte("Error: error when parsing the list"))
		return
	}

	old := LoadList(l.Id)
	if old == nil {
		w.Write([]byte("Error: this list doesn't exists"))
		return
	}
	if l.ModificationTime < old.ModificationTime {
		w.Write([]byte("Error: modification conflict"))
		return
	}

	l.Save()
	w.Write([]byte("Success"))
}

func main() {
	InitModel()

	http.HandleFunc("/", Index)
	http.HandleFunc("/about", About)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/view/", View)
	http.HandleFunc("/save", Save)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(staticDir))))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Cannot start HTTP server:", err)
	}
}
