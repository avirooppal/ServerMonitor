package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

// ... (omitted types)

func InitDB() {
	os.MkdirAll("data", 0755)
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
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createSystemsTable := `CREATE TABLE IF NOT EXISTS systems (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
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

	InitSessionsTable()
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

// System Management

func AddSystem(userID int, name, url, apiKey string) (int64, error) {
	res, err := DB.Exec("INSERT INTO systems (user_id, name, url, api_key) VALUES (?, ?, ?, ?)", userID, name, url, apiKey)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetSystems(userID int) ([]System, error) {
	rows, err := DB.Query("SELECT id, user_id, name, url, api_key, created_at FROM systems WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var systems []System
	for rows.Next() {
		var s System
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.URL, &s.APIKey, &s.CreatedAt); err != nil {
			return nil, err
		}
		systems = append(systems, s)
	}
	return systems, nil
}

func GetSystem(id int) (*System, error) {
	var s System
	err := DB.QueryRow("SELECT id, user_id, name, url, api_key, created_at FROM systems WHERE id = ?", id).Scan(&s.ID, &s.UserID, &s.Name, &s.URL, &s.APIKey, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetSystemByAPIKey(apiKey string) (*System, error) {
	var s System
	err := DB.QueryRow("SELECT id, user_id, name, url, api_key, created_at FROM systems WHERE api_key = ?", apiKey).Scan(&s.ID, &s.UserID, &s.Name, &s.URL, &s.APIKey, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteSystem(id, userID int) error {
	_, err := DB.Exec("DELETE FROM systems WHERE id = ? AND user_id = ?", id, userID)
	return err
}

// User Management

func CreateUser(email, passwordHash string) error {
	_, err := DB.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, passwordHash)
	return err
}

func GetUserByEmail(email string) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT id, email, password_hash, created_at FROM users WHERE email = ?", email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByID(id int) (*User, error) {
	var u User
	err := DB.QueryRow("SELECT id, email, password_hash, created_at FROM users WHERE id = ?", id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Session Management

type Session struct {
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func InitSessionsTable() {
	createTable := `CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
	if _, err := DB.Exec(createTable); err != nil {
		log.Fatalf("Failed to create sessions table: %v", err)
	}
}

func CreateSession(token string, userID int, expiresAt time.Time) error {
	_, err := DB.Exec("INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)", token, userID, expiresAt)
	return err
}

func GetSession(token string) (*Session, error) {
	var s Session
	err := DB.QueryRow("SELECT token, user_id, expires_at FROM sessions WHERE token = ?", token).Scan(&s.Token, &s.UserID, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteSession(token string) error {
	_, err := DB.Exec("DELETE FROM sessions WHERE token = ?", token)
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
