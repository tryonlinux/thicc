package database

const schema = `
CREATE TABLE IF NOT EXISTS weights (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    weight REAL NOT NULL,
    bmi REAL NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_weights_date ON weights(date DESC);

CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

// InitializeSchema creates all tables
func InitializeSchema(db *DB) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	return nil
}
