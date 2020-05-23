package main

import (
	"checklists/aircraft"
	"log"
	"net/http"
)

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

const basePath = "/api"

func main() {
	aircraft.SetupRoutes(basePath)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
