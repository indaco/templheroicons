package main

import (
	"log"
	"net/http"
	"os"

	"github.com/indaco/templheroicons/_demos/pages"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	err := pages.HomePage().Render(r.Context(), w)
	if err != nil {
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", HandleHome)

	port := ":3300"
	log.Printf("Listening on %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Printf("failed to start server: %v", err)
		os.Exit(1)
	}
}
