package db

import (
	"database/sql"
	"log"


	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	dbPath := "server-moni.db"
	// Check if running in docker or specific path needed, for now relative to binary
	
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
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
