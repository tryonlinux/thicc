package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// DB wraps the sql.DB connection
type DB struct {
	*sql.DB
}

// Open opens a connection to the SQLite database and initializes the schema
func Open(dbPath string) (*DB, error) {
	sqlDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	db := &DB{sqlDB}

	// Initialize schema
	if err := InitializeSchema(db); err != nil {
		sqlDB.Close()
		return nil, err
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
