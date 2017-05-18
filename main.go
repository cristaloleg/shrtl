package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"sync"

	"github.com/pressly/chi"
)

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
}

func submit(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	name := db.Add(url)
	w.Write([]byte(name))
}

func open(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]

	if url, ok := db.Get(name); ok {
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("URL Not Found"))
}

var db storage

type storage struct {
	sync.RWMutex
	data map[string]string
}

func (s *storage) Add(url string) (name string) {
	name = generateName()
	s.Lock()
	s.data[name] = url
	s.Unlock()
	return name
}

func (s *storage) Get(name string) (url string, ok bool) {
	s.RLock()
	url, ok = s.data[name]
	s.RUnlock()
	return url, ok
}

func generateName() string {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return string(b)
}
