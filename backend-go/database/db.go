package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	// Create users table
	db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		name TEXT NOT NULL,
		role TEXT NOT NULL CHECK(role IN ('worker', 'admin'))
	);`)

	// Create shifts table
	db.Exec(`
	CREATE TABLE IF NOT EXISTS shifts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATE NOT NULL,
		start TIME NOT NULL,
		end TIME NOT NULL,
		role TEXT NOT NULL,
		location TEXT
	);`)

	// Create requests table
	db.Exec(`
	CREATE TABLE IF NOT EXISTS requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		shift_id INTEGER NOT NULL,
		status TEXT NOT NULL CHECK(status IN ('pending', 'approved', 'rejected')),
		FOREIGN KEY(user_id) REFERENCES users(id),
		FOREIGN KEY(shift_id) REFERENCES shifts(id)
	);`)

	return db
}
