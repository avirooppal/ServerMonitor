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

	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		api_key TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// We need to add user_id to systems. 
	// For existing tables, we might need a migration, but for this "clean slate" approach:
	createSystemsTable := `CREATE TABLE IF NOT EXISTS systems (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		api_key TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	if _, err := DB.Exec(createConfigTable); err != nil {
		log.Fatalf("Failed to create config table: %v", err)
	}

	if _, err := DB.Exec(createUsersTable); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	if _, err := DB.Exec(createSystemsTable); err != nil {
		log.Fatalf("Failed to create systems table: %v", err)
	}

	// Migration: Check if user_id column exists in systems (for existing deployments)
	// Simple check: try to query it. If fails, add it.
	_, err = DB.Query("SELECT user_id FROM systems LIMIT 1")
	if err != nil {
		// Column likely missing
		log.Println("Migrating systems table: adding user_id column...")
		_, err = DB.Exec("ALTER TABLE systems ADD COLUMN user_id INTEGER REFERENCES users(id)")
		if err != nil {
			log.Printf("Warning: Failed to add user_id column (might already exist or other error): %v", err)
		}
	}

	InitDiskHistoryTable()
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

func AddSystem(userID int, name, url, apiKey string) (int64, error) {
	res, err := DB.Exec("INSERT INTO systems (user_id, name, url, api_key) VALUES (?, ?, ?, ?)", userID, name, url, apiKey)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetSystems(userID int) ([]System, error) {
	rows, err := DB.Query("SELECT id, name, url, api_key, created_at FROM systems WHERE user_id = ? ORDER BY created_at DESC", userID)
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

func GetSystem(userID int, id string) (*System, error) {
	var s System
	// id can be string from URL param, let's assume caller converts or we use string in query
	err := DB.QueryRow("SELECT id, name, url, api_key, created_at FROM systems WHERE id = ? AND user_id = ?", id, userID).Scan(&s.ID, &s.Name, &s.URL, &s.APIKey, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetSystemByAPIKey(apiKey string) (*System, error) {
	var s System
	err := DB.QueryRow("SELECT id, name, url, api_key, created_at FROM systems WHERE api_key = ?", apiKey).Scan(&s.ID, &s.Name, &s.URL, &s.APIKey, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteSystem(userID int, id string) error {
	_, err := DB.Exec("DELETE FROM systems WHERE id = ? AND user_id = ?", id, userID)
	return err
}

// Disk History

type DiskHistory struct {
	ID          int64     `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	UsedPercent float64   `json:"used_percent"`
	Total       uint64    `json:"total"`
	Used        uint64    `json:"used"`
}

func InitDiskHistoryTable() {
	createTable := `CREATE TABLE IF NOT EXISTS disk_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		used_percent REAL,
		total INTEGER,
		used INTEGER
	);`
	if _, err := DB.Exec(createTable); err != nil {
		log.Printf("Failed to create disk_history table: %v", err)
	}
}

func AddDiskHistory(usedPercent float64, total, used uint64) error {
	_, err := DB.Exec("INSERT INTO disk_history (used_percent, total, used) VALUES (?, ?, ?)", usedPercent, total, used)
	return err
}

func GetDiskHistory(limit int) ([]DiskHistory, error) {
	rows, err := DB.Query("SELECT id, timestamp, used_percent, total, used FROM disk_history ORDER BY timestamp DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []DiskHistory
	for rows.Next() {
		var h DiskHistory
		if err := rows.Scan(&h.ID, &h.Timestamp, &h.UsedPercent, &h.Total, &h.Used); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}
