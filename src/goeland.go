package main

import (
	"fmt"
	"bytes"
	"net/http"
	"html/template"
)

const templatesDir = "../templates"
const staticDir = "../static"

type Page struct {
	Content template.HTML
}

func renderTemplate(w http.ResponseWriter, tmpl string, arg interface{}) {
	skel, _ := template.ParseFiles(fmt.Sprintf("%s/skeleton.html", templatesDir))
	content, _ := template.ParseFiles(fmt.Sprintf("%s/%s.html", templatesDir, tmpl));
	
	buff := &bytes.Buffer{}
	content.Execute(buff, arg)
	c := template.HTML(buff.String())
	page := &Page{Content: c}
	skel.Execute(w, page)
}

func index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func about(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	l := newList(name)
	l.save()
	http.Redirect(w, r, "/view/" + l.Id, http.StatusFound)
}

func view(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[6:]
	l := loadList(id)
	if (l != nil) {
		renderTemplate(w, "view", l)
	} else {
		renderTemplate(w, "message",
			&Message{"This list doesn't exists", "Back", "/"})
	}
}

func save(w http.ResponseWriter, r *http.Request) {
	l := parseList([]byte(r.FormValue("list")))
	if l == nil {
		w.Write([]byte("Error: error when parsing the list"))
		return
	}

	old := loadList(l.Id)
	if old == nil {
		w.Write([]byte("Error: this list doesn't exists"))
		return
	}
	if l.ModificationTime < old.ModificationTime {
		w.Write([]byte("Error: modification conflict"))
		return
	}

	l.save()
	w.Write([]byte("Success"))
}

func main() {
	init()

	http.HandleFunc("/", index)
	http.HandleFunc("/about", about)
	http.HandleFunc("/create", create)
	http.HandleFunc("/view/", view)
	http.HandleFunc("/save", save)
	
	http.Handle("/static/",
		http.StripPrefix("/static/",
		http.FileServer(http.Dir(staticDir))))
	http.ListenAndServe(":8080", nil)
}
