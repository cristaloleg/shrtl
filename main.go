package main

import (
	"crypto/md5"
	"hash"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pressly/chi"
)

var db storage
var h hash.Hash
var stats map[string]uint64
var ch chan string

func init() {
	db = *newStorage()
	h = md5.New()
	stats = make(map[string]uint64)
	ch = make(chan string, 1024)

	go func() {
		for {
			name := <-ch
			value, ok := stats[name]
			if ok {
				value++
				stats[name] = value
			}
		}
	}()
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
	w.Write([]byte(name))
}

func open(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]

	if url, ok := db.Get(name); ok {
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
		statsUpdate(name)
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("URL Not Found"))
}

type storage struct {
	sync.RWMutex
	data map[string]string
}

func newStorage() *storage {
	return &storage{
		data: make(map[string]string),
	}
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
	cputime := time.Now().Unix()
	return strconv.FormatInt(cputime, 36)
}

func statsUpdate(name string) {
	ch <- name
}
