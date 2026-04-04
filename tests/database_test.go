package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tryonlinux/thicc/internal/database"
)

func TestOpenError(t *testing.T) {
	// Provide an invalid path where it's impossible to create a file
	// e.g., a path that doesn't exist and can't be created
	invalidPath := "/nonexistent/path/weights.db"
	db, err := database.Open(invalidPath)
	if err == nil {
		db.Close()
		t.Errorf("Open(%q) should have returned an error", invalidPath)
	}
}

func TestInitializeSchema(t *testing.T) {
	// Create a temporary database file
	tmpDir, err := os.MkdirTemp("", "thicc_db_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Verify tables exist
	tables := []string{"weights", "settings"}
	for _, table := range tables {
		var name string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("Table %s does not exist: %v", table, err)
		}
	}
}
