package db

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

type System struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

func InitDB() {
	dbPath := "data/server-moni.db"
	
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createConfigTable := `CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT
	);`

	createSystemsTable := `CREATE TABLE IF NOT EXISTS systems (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		api_key TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(createConfigTable); err != nil {
		log.Fatalf("Failed to create config table: %v", err)
	}

	if _, err := DB.Exec(createSystemsTable); err != nil {
		log.Fatalf("Failed to create systems table: %v", err)
	}
}

func GetConfig(key string) (string, error) {
	var value string
	err := DB.QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
	return value, err
}

func SetConfig(key, value string) error {
	_, err := DB.Exec("INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)", key, value)
	return err
}

func AddSystem(name, url, apiKey string) (int64, error) {
	res, err := DB.Exec("INSERT INTO systems (name, url, api_key) VALUES (?, ?, ?)", name, url, apiKey)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetSystems() ([]System, error) {
	rows, err := DB.Query("SELECT id, name, url, api_key, created_at FROM systems ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var systems []System
	for rows.Next() {
		var s System
		if err := rows.Scan(&s.ID, &s.Name, &s.URL, &s.APIKey, &s.CreatedAt); err != nil {
			return nil, err
		}
		systems = append(systems, s)
	}
	return systems, nil
}

func GetSystem(id int) (*System, error) {
	var s System
	err := DB.QueryRow("SELECT id, name, url, api_key, created_at FROM systems WHERE id = ?", id).Scan(&s.ID, &s.Name, &s.URL, &s.APIKey, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteSystem(id int) error {
	_, err := DB.Exec("DELETE FROM systems WHERE id = ?", id)
	return err
}
