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
	handleChecklists := http.HandlerFunc(checklistsHandler)
	handleChecklistForAircraft := http.HandlerFunc(checklistHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiPath, aircraftBaseRoutePath), cors.Middleware(handleAllAircraft))
	http.Handle(fmt.Sprintf("%s/%s/", apiPath, aircraftBaseRoutePath), cors.Middleware(handleAircraft))
	http.Handle(fmt.Sprintf("%s/%s", apiPath, "checklists"), cors.Middleware(handleChecklists))
	http.Handle(fmt.Sprintf("%s/%s/", apiPath, "checklists"), cors.Middleware(handleChecklistForAircraft))
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
	urlPath := strings.Split(r.URL.Path, "/")
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

	// Switch for CRUD
	switch r.Method {
	case http.MethodGet:
		// aircraft.Checklists = []string{fakeChecklist.Title}
		checklists, err := getChecklistsForAircraft(aircraft.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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

func checklistHandler(w http.ResponseWriter, r *http.Request) {
	checklistPath := strings.Split(r.URL.Path, "checklists/")
	checklistID, err := strconv.Atoi(checklistPath[len(checklistPath)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		checklistData, err := getChecklistDetailByID(checklistID)
		if checklistData == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		ckJSON, err := json.Marshal(checklistData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(ckJSON)
		return
	case http.MethodPut:
		checklistData, err := getChecklistDetailByID(checklistID)
		if checklistData == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var updatedChecklist Checklist
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(requestBody, &updatedChecklist)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateChecklist(checklistID, updatedChecklist)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

}

func checklistsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var newChecklist Checklist
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(requestBody, &newChecklist)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = addChecklist(newChecklist)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
}
