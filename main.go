package main

import (
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
}

func open(w http.ResponseWriter, r *http.Request) {
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
	return ""
}
