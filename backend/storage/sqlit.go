package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"kernelscope/models"

	// Keep the pure Go driver we set up earlier
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDatabase() error {
	var err error

	// Initializing with the pure Go driver tag
	DB, err = sql.Open("sqlite", "kernelscope.db")
	if err != nil {
		return err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS timeline_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp INTEGER,
		type TEXT,
		pid INTEGER,
		process TEXT,
		importance TEXT,
		details TEXT
	);
	`

	_, err = DB.Exec(createTable)
	if err != nil {
		return err
	}

	fmt.Println("✅ SQLite initialized (Pure Go)")

	return nil
}

func SaveEvent(event models.Event) error {
	detailsJSON, err := json.Marshal(event.Details)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO timeline_events 
	(timestamp, type, pid, process, importance, details) 
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err = DB.Exec(
		query,
		event.Timestamp,
		event.Type,
		event.PID,
		event.Process,
		event.Importance,
		string(detailsJSON),
	)

	return err
}
func LoadEvents(limit int) ([]models.Event, error) {
	query := `
	SELECT timestamp, type, pid, process, importance, details
	FROM timeline_events
	ORDER BY id DESC
	LIMIT ?
	`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event
		var eventType string
		var detailsJSON string

		err := rows.Scan(
			&event.Timestamp,
			&eventType,
			&event.PID,
			&event.Process,
			&event.Importance,
			&detailsJSON,
		)
		if err != nil {
			return nil, err
		}

		event.Type = models.EventType(eventType)

		if detailsJSON != "" {
			_ = json.Unmarshal([]byte(detailsJSON), &event.Details)
		}

		events = append(events, event)
	}

	return events, rows.Err()
}