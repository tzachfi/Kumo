// Command server is the HTTP entrypoint for the Kumo SDUI API.
package main

import (
	"log"
	"net/http"

	"github.com/tzachfi/kumo/server/internal/api/handler"
	"github.com/tzachfi/kumo/server/internal/store"
)

func main() {
	mux := http.NewServeMux()

	st := store.NewMock()
	jh := &handler.JourneyHandler{Store: st}

	mux.HandleFunc("GET /api/journey/{id}", jh.Get)
	mux.HandleFunc("POST /api/journey/generate", jh.Generate)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
