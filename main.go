package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/pressly/chi"
)

var db storage
var st stats

func init() {
	db = newStorage()
	st = newStats()
}

func main() {
	r := chi.NewRouter()

	r.Get("/", home)
	r.Post("/", submit)
	r.Get("/:name", open)
	r.Post("/:name", open)

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, nil)
}

func submit(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	name := db.Add(url)
	st.Add(url)
	w.Write([]byte(name))
}

func open(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]

	if url, ok := db.Get(name); ok {
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
		st.Inc(name)
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("URL Not Found"))
}
