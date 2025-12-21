package config

import (
	"os"
	"path/filepath"
)

// GetDatabasePath returns the path to the SQLite database file.
// Creates the ~/.thicc directory if it doesn't exist.
func GetDatabasePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	thiccDir := filepath.Join(homeDir, ".thicc")

	// Create .thicc directory if it doesn't exist (0700 = owner only for security)
	if err := os.MkdirAll(thiccDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(thiccDir, "weights.db"), nil
}
