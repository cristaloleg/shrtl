package main

import (
	"log"
	"net/http"

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
