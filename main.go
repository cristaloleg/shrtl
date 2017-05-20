package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var db storage
var st stats
var homeTmpl *template.Template

func init() {
	db = newStorage()
	st = newStats()
	homeTmpl, _ = template.ParseFiles("home.html")
}

func main() {
	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	method := r.Method

	if url == "/" {
		switch {
		case method == "GET":
			homeTmpl.Execute(w, nil)

		case method == "POST":
			url := r.PostFormValue("url")
			if !strings.HasPrefix("http://", url) {
				url = "http://" + url
			}
			name := db.Add(url)
			st.Add(url)
			w.Write([]byte(name))

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		if method != "GET" && method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		name := url[1:]

		if url, ok := db.Get(name); ok {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
			st.Inc(name)
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("URL Not Found"))
	}
}
