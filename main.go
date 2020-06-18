package main

import (
	"checklists/aircraft"
	"checklists/database"
	"fmt"
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
	const PORT = "8080"
	const ADDR = "localhost:" + PORT
	fmt.Printf("Server started on port: %s\n", PORT)
	log.Fatal(http.ListenAndServe(ADDR, nil))
}
