package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Aircraft struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"alias"`
}

// First, let's do CRUD on the Aircraft. For now, it's just an id, name, and short name.
var allAircraft []Aircraft

func getAircraftByID(aircraftID int) (*Aircraft, int) {
	for i, aircraft := range allAircraft {
		if aircraft.ID == aircraftID {
			return &aircraft, i
		}
	}
	return nil, 0
}

// Marshal: struct -> JSON
// Unmarshal: JSON -> struct (into a byte slice)
func allAircraftHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		acJSON, err := json.Marshal(allAircraft)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(acJSON)
	case http.MethodPost:
		var newAircraft Aircraft
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(reqBody, &newAircraft)
		allAircraft = append(allAircraft, newAircraft)
		fmt.Println(allAircraft)
		w.WriteHeader(http.StatusCreated)
	}
}

func singleAircraftHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Split(r.URL.Path, "aircraft/")
	aircraftID, err := strconv.Atoi(urlPath[len(urlPath)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	aircraft, acIndex := getAircraftByID(aircraftID)
	if aircraft == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Switch for CRUD
	switch r.Method {
	case http.MethodGet:
		fmt.Println("Getting aircraft", aircraftID)
		aircraftJSON, err := json.Marshal(aircraft)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(aircraftJSON)
		return
	case http.MethodPut:
		var changedAC Aircraft
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &changedAC)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if changedAC.ID != aircraftID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		aircraft = &changedAC
		allAircraft[acIndex] = *aircraft
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		w.WriteHeader(http.StatusNotImplemented)
		return
	default:
		w.WriteHeader(http.StatusTeapot)
	}
}

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

func main() {
	fname := "aircraft.json"
	acFromFile, _ := ioutil.ReadFile(fname)
	err := json.Unmarshal(acFromFile, &allAircraft)
	if err != nil {
		log.Fatal("Error reading JSON")
	}
	allHandler := http.HandlerFunc(allAircraftHandler)
	singleHandler := http.HandlerFunc(singleAircraftHandler)
	http.Handle("/aircraft", middlewareHandler(allHandler))
	http.Handle("/aircraft/", middlewareHandler(singleHandler))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
