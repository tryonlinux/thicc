package tests

import (
	"strings"
	"testing"

	"github.com/tryonlinux/thicc/internal/config"
)

func TestGetDatabasePath(t *testing.T) {
	path, err := config.GetDatabasePath()
	if err != nil {
		t.Fatalf("GetDatabasePath() error = %v", err)
	}

	if path == "" {
		t.Fatal("GetDatabasePath() returned empty string")
	}

	// The path should end with .thicc/weights.db (on Unix) or .thicc\weights.db (on Windows)
	expectedSuffix := ".thicc/weights.db"
	expectedSuffixWin := ".thicc\\weights.db"

	if !strings.HasSuffix(path, expectedSuffix) && !strings.HasSuffix(path, expectedSuffixWin) {
		t.Errorf("GetDatabasePath() = %q, want it to end with %q", path, expectedSuffix)
	}
}
