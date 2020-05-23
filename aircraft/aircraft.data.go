package aircraft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

var aircraftMap = struct {
	sync.RWMutex
	m map[int]Aircraft
}{m: make(map[int]Aircraft)}

func init() {
	fmt.Println("Loading aircraft..")
	acMap, err := loadAircraftMap()
	aircraftMap.m = acMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d aircraft loaded...", len(aircraftMap.m))
}

func loadAircraftMap() (map[int]Aircraft, error) {
	fileName := "aircraft.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("File [%s] does not exist", fileName)
	}

	file, _ := ioutil.ReadFile(fileName)
	aircraftList := make([]Aircraft, 0)
	err = json.Unmarshal([]byte(file), &aircraftList)
	if err != nil {
		log.Fatal(err)
	}
	aircraftMap := make(map[int]Aircraft)
	for i := 0; i < len(aircraftList); i++ {
		aircraftMap[aircraftList[i].ID] = aircraftList[i]
	}
	return aircraftMap, nil
}

func getAircraftByID(aircraftID int) *Aircraft {
	aircraftMap.RLock()
	defer aircraftMap.RUnlock()
	if aircraft, ok := aircraftMap.m[aircraftID]; ok {
		return &aircraft
	}
	return nil
}

func removeAircraft(aircraftID int) {
	aircraftMap.Lock()
	defer aircraftMap.Unlock()
	delete(aircraftMap.m, aircraftID)
}

func getAircraftList() []Aircraft {
	aircraftMap.RLock()
	aircraft := make([]Aircraft, 0, len(aircraftMap.m))
	for _, value := range aircraftMap.m {
		aircraft = append(aircraft, value)
	}
	aircraftMap.RUnlock()
	return aircraft
}

func getAircraftIDs() []int {
	aircraftMap.RLock()
	aircraftIDs := []int{}
	for key := range aircraftMap.m {
		aircraftIDs = append(aircraftIDs, key)
	}
	aircraftMap.Unlock()
	sort.Ints(aircraftIDs)
	return aircraftIDs
}

func getNextAircraftID() int {
	aircraftIDs := getAircraftIDs()
	return aircraftIDs[len(aircraftIDs)-1] + 1
}

func addOrUpdateAircraft(aircraft Aircraft) (int, error) {
	aircraftID := -1
	if aircraft.ID > 0 {
		oldAircraft := getAircraftByID(aircraft.ID)
		if oldAircraft == nil {
			return 0, fmt.Errorf("Aircraft [%d] does not exist", aircraft.ID)
		}
		aircraftID = oldAircraft.ID
	} else {
		aircraftID = getNextAircraftID()
		aircraft.ID = aircraftID
	}
	aircraftMap.Lock()
	aircraftMap.m[aircraftID] = aircraft
	aircraftMap.Unlock()
	return aircraftID, nil
}
