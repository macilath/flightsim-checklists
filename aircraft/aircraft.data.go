package aircraft

import (
	"checklists/database"
	"context"
	"database/sql"
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
		allAircraft = append(allAircraft, aircraft)
	}
	return allAircraft, nil
}

func addAircraft(aircraft Aircraft) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var retId int // Assumption that this is the most recent Id - if col is not set to autoincrement it may not be
	err := database.DbConn.QueryRowContext(ctx, `INSERT INTO aircraft (id, name, alias) VALUES ($1, $2, $3) RETURNING id`, aircraft.ID, aircraft.Name, aircraft.ShortName).Scan(&retId)
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, nil
	}
	return int(retId), nil
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
