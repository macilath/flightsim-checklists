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
	//handleChecklistsForAircraft := http.HandlerFunc(checklistsHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiPath, aircraftBaseRoutePath), cors.Middleware(handleAllAircraft))
	http.Handle(fmt.Sprintf("%s/%s/", apiPath, aircraftBaseRoutePath), cors.Middleware(handleAircraft))
	//http.Handle(fmt.Sprintf("%s/%s/checklists", apiPath, aircraftBaseRoutePath, cors.Middleware(handleChecklistsForAircraft)))
}

func allAircraftHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		allAircraft, err := getAircraftList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
		aircraftID, err := addAircraft(newAircraft)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"aircraftId":%d}`, aircraftID)))
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func singleAircraftHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Split(r.URL.Path, "aircraft/")
	aircraftID, err := strconv.Atoi(urlPath[len(urlPath)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// First load the aircraft metadata
	aircraft, err := getAircraftByID(aircraftID)
	if aircraft == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Load what checklist(s) are available - just the names, though
	// items := []string{
	// 	"PARKING BRAKE - SET",
	// 	"ANTI-SKID - ON",
	// 	"MASTER ARM - SAFE",
	// 	"WING FOLD - MATCH",
	// 	"BATTERY - ON",
	// 	"BRAKE PRESSURE - 3000",
	// 	"FIRE TEST A - PERFORM",
	// 	"FIRE TEST B - PERFORM",
	// 	"APU - START"}
	// fakeChecklist := Checklist{Title: "PRE-START CHECKLIST", Items: items}

	// Switch for CRUD
	switch r.Method {
	case http.MethodGet:
		// aircraft.Checklists = []string{fakeChecklist.Title}
		checklists, err := getChecklistsForAircraft(aircraft.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		checkMe, err := getChecklistDetailByIDs(0, 1)
		fmt.Println(len(checkMe.Items))
		aircraft.Checklists = checklists
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
		updateAircraft(changedAC)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removeAircraft(aircraftID)
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusTeapot)
	}
}

// func checklistsHandler(w http.ResponseWriter, r *http.Request) {

// }
