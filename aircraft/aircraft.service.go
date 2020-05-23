package aircraft

import (
	"checklists/cors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const aircraftBaseRoutePath = "aircraft"

func SetupRoutes(apiPath string) {
	handleAllAircraft := http.HandlerFunc(allAircraftHandler)
	handleAircraft := http.HandlerFunc(singleAircraftHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiPath, aircraftBaseRoutePath), cors.Middleware(handleAllAircraft))
	http.Handle(fmt.Sprintf("%s/%s/", apiPath, aircraftBaseRoutePath), cors.Middleware(handleAircraft))
}

func allAircraftHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		allAircraft := getAircraftList()
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
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newAircraft.ID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = addOrUpdateAircraft(newAircraft)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	}
}

func singleAircraftHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Split(r.URL.Path, "aircraft/")
	aircraftID, err := strconv.Atoi(urlPath[len(urlPath)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	aircraft := getAircraftByID(aircraftID)
	if aircraft == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Switch for CRUD
	switch r.Method {
	case http.MethodGet:
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
		addOrUpdateAircraft(changedAC)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removeAircraft(aircraftID)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusTeapot)
	}
}
