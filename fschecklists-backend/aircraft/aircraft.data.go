package aircraft

import (
	"checklists/database"
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

func getAircraftByID(aircraftID int) (*Aircraft, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result := database.DbConn.QueryRowContext(ctx, `SELECT id, name, alias FROM aircraft WHERE id = $1`, aircraftID)
	aircraft := &Aircraft{}
	err := result.Scan(&aircraft.ID, &aircraft.Name, &aircraft.ShortName)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return aircraft, nil
}

func removeAircraft(aircraftID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM aircraft WHERE id = $1`, aircraftID)
	if err != nil {
		return err
	}
	return nil
}

func getAircraftList() ([]Aircraft, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT id, name, alias FROM aircraft`)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	allAircraft := make([]Aircraft, 0)
	for results.Next() {
		var aircraft Aircraft
		results.Scan(&aircraft.ID, &aircraft.Name, &aircraft.ShortName)
		cls, err := getChecklistsForAircraft(aircraft.ID)
		if err != nil {
			cls = make([]LiteChecklist, 0)
		} else if cls == nil {
			cls = make([]LiteChecklist, 0)
		}
		aircraft.Checklists = cls
		allAircraft = append(allAircraft, aircraft)
	}
	return allAircraft, nil
}

func addAircraft(aircraft Aircraft) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var retID int // Assumption that this is the most recent Id - if col is not set to autoincrement it may not be
	err := database.DbConn.QueryRowContext(ctx,
		`INSERT INTO aircraft (id, name, alias) VALUES ($1, $2, $3) RETURNING id`,
		aircraft.ID,
		aircraft.Name,
		aircraft.ShortName,
	).Scan(&retID)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, nil
	}
	return int(retID), nil
}

func updateAircraft(aircraft Aircraft) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `UPDATE aircraft SET id = $1, name = $2, alias = $3 WHERE id = $1`, aircraft.ID, aircraft.Name, aircraft.ShortName)
	if err != nil {
		return err
	}
	return nil
}

func getChecklistsForAircraft(aircraftID int) ([]LiteChecklist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT title, id FROM checklist WHERE aircraft_id = $1`, aircraftID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer results.Close()
	checklists := make([]LiteChecklist, 0)
	for results.Next() {
		var clist LiteChecklist
		results.Scan(&clist.Title, &clist.ID)
		checklists = append(checklists, clist)
	}
	return checklists, nil
}

func getChecklistDetailByID(checklistID int) (*Checklist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result := database.DbConn.QueryRowContext(ctx, `SELECT id, title, items, aircraft_id FROM checklist WHERE id = $1`, checklistID)

	type dbChecklist struct {
		Items string // Comes as a stringified JSON
	}

	checkListDetail := Checklist{}
	checklist := dbChecklist{}
	err := result.Scan(&checkListDetail.ID, &checkListDetail.Title, &checklist.Items, &checkListDetail.AircraftID)
	if err != nil {
		return &Checklist{}, err
	}
	err = json.Unmarshal([]byte(checklist.Items), &checkListDetail.Items)

	if err != nil {
		return &Checklist{}, err
	}
	return &checkListDetail, nil
}

func addChecklist(checklist Checklist) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// pack items into its own object
	items, err := json.Marshal(checklist.Items)

	var retID int
	err = database.DbConn.QueryRowContext(ctx, `INSERT INTO checklist (title, items, aircraft_id) VALUES ($1, $2, $3) RETURNING id`, checklist.Title, items, checklist.AircraftID).Scan(&retID)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, nil
	}
	return int(retID), nil
}

func updateChecklist(id int, content Checklist) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	items, err := json.Marshal(content.Items)
	if err != nil {
		return err
	}
	_, err = database.DbConn.ExecContext(ctx, `UPDATE checklist SET title = $1, items=$2, aircraft_id=$3 WHERE id=$4`, content.Title, items, content.AircraftID, id)

	return err
}
