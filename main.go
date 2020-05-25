package main

import (
	"checklists/aircraft"
	"checklists/database"
	"log"
	"net/http"

	_ "github.com/lib/pq" // importing for side effects
)

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

const basePath = "/api"

func main() {
	database.SetupDatabase()
	aircraft.SetupRoutes(basePath)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
