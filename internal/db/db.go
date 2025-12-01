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
		email TEXT UNIQUE,
		password_hash TEXT,
		provider TEXT DEFAULT 'local',
		provider_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

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
	
	// Migration: Add user_id to systems if it doesn't exist
	// SQLite doesn't support IF NOT EXISTS for columns, so we ignore error
	_, _ = DB.Exec("ALTER TABLE systems ADD COLUMN user_id INTEGER REFERENCES users(id)")

	InitDiskHistoryTable()
}

// User Struct
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Provider     string    `json:"provider"`
	ProviderID   string    `json:"provider_id"`
	CreatedAt    time.Time `json:"created_at"`
}

func CreateUser(email, passwordHash, provider, providerID string) (int64, error) {
	res, err := DB.Exec("INSERT INTO users (email, password_hash, provider, provider_id) VALUES (?, ?, ?, ?)", email, passwordHash, provider, providerID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetUserByEmail(email string) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT id, email, password_hash, provider, provider_id, created_at FROM users WHERE email = ?", email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Provider, &u.ProviderID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByID(id int) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT id, email, password_hash, provider, provider_id, created_at FROM users WHERE id = ?", id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Provider, &u.ProviderID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
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
	// If userID is 0 (Self-Hosted mode legacy), we might want to return all or handle differently.
	// But for now, we assume strict separation if userID is provided.
	// For backward compatibility/Self-Hosted, we might pass 0 or -1.
	
	query := "SELECT id, name, url, api_key, created_at FROM systems WHERE user_id = ? ORDER BY created_at DESC"
	args := []interface{}{userID}
	
	if userID == 0 {
		// Fetch all (Self-Hosted Mode) or systems with NULL user_id
		query = "SELECT id, name, url, api_key, created_at FROM systems WHERE user_id IS NULL OR user_id = 0 ORDER BY created_at DESC"
		args = []interface{}{}
	}

	rows, err := DB.Query(query, args...)
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

func GetSystem(id int, userID int) (*System, error) {
	var s System
	query := "SELECT id, name, url, api_key, created_at FROM systems WHERE id = ? AND user_id = ?"
	args := []interface{}{id, userID}

	if userID == 0 {
		query = "SELECT id, name, url, api_key, created_at FROM systems WHERE id = ?"
		args = []interface{}{id}
	}

	err := DB.QueryRow(query, args...).Scan(&s.ID, &s.Name, &s.URL, &s.APIKey, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteSystem(id int, userID int) error {
	query := "DELETE FROM systems WHERE id = ? AND user_id = ?"
	args := []interface{}{id, userID}
	
	if userID == 0 {
		query = "DELETE FROM systems WHERE id = ?"
		args = []interface{}{id}
	}
	
	_, err := DB.Exec(query, args...)
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
