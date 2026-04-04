package tests

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/tryonlinux/thicc/cmd"
)

func TestCLICommands(t *testing.T) {
	// Setup temporary home directory
	tmpHome, err := os.MkdirTemp("", "thicc_test_home_*")
	if err != nil {
		t.Fatalf("Failed to create temp home: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	// Save original environment and restore after test
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", oldHome)

	// Capture stdout
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Simulate first-time setup input
	stdinR, stdinW, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = stdinR
	go func() {
		stdinW.Write([]byte("lbs\nin\n70\n150\n"))
		stdinW.Close()
	}()

	// Run root command which triggers initialization and setup
	os.Args = []string{"thicc"}
	cmd.Execute()

	w.Close()
	os.Stdout = stdout
	os.Stdin = oldStdin

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "First Time Setup") {
		t.Error("Expected setup prompt, got:", output)
	}

	// Now test 'add' command
	stdout = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w

	os.Args = []string{"thicc", "add", "160.5"}
	cmd.Execute()

	w.Close()
	os.Stdout = stdout
	buf.Reset()
	io.Copy(&buf, r)
	output = buf.String()

	if !strings.Contains(output, "Added weight: 160.50 lbs") {
		t.Error("Expected add confirmation, got:", output)
	}

	// Test 'add' with date
	stdout = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w

	os.Args = []string{"thicc", "add", "158.0", "2024-01-01"}
	cmd.Execute()

	w.Close()
	os.Stdout = stdout
	buf.Reset()
	io.Copy(&buf, r)
	output = buf.String()
	if !strings.Contains(output, "Added weight: 158.00 lbs on 2024-01-01") {
		t.Error("Expected add confirmation with date, got:", output)
	}

	// Test 'add' with invalid weight
	os.Args = []string{"thicc", "add", "--", "-1"}
	cmd.Execute()
	os.Args = []string{"thicc", "add", "0"}
	cmd.Execute()
	os.Args = []string{"thicc", "add", "2000"}
	cmd.Execute()

	// Test 'add' with invalid date
	os.Args = []string{"thicc", "add", "160", "invalid-date"}
	cmd.Execute()

	// Test 'show'
	stdout = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w

	os.Args = []string{"thicc", "show"}
	cmd.Execute()

	w.Close()
	os.Stdout = stdout
	buf.Reset()
	io.Copy(&buf, r)
	output = buf.String()
	if !strings.Contains(output, "Weight Tracker") {
		t.Error("Expected Weight Tracker output, got:", output)
	}

	// Test 'show' with limit
	os.Args = []string{"thicc", "show", "10"}
	cmd.Execute()
	os.Args = []string{"thicc", "show", "0"}
	cmd.Execute()
	os.Args = []string{"thicc", "show", "--", "-1"}
	cmd.Execute()

	// Test 'show' with date
	os.Args = []string{"thicc", "show", "2024-01-01"}
	cmd.Execute()
	os.Args = []string{"thicc", "show", "invalid-arg"}
	cmd.Execute()

	// Test 'goal'
	stdout = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w

	os.Args = []string{"thicc", "goal", "145"}
	cmd.Execute()

	w.Close()
	os.Stdout = stdout
	buf.Reset()
	io.Copy(&buf, r)
	output = buf.String()
	if !strings.Contains(output, "Goal weight set to 145.00 lbs") {
		t.Error("Expected goal set confirmation, got:", output)
	}

	// Test 'goal' with invalid weight
	os.Args = []string{"thicc", "goal", "0"}
	cmd.Execute()

	// Test 'modify'
	// ID 1 should be the first weight we added
	os.Args = []string{"thicc", "modify", "1", "159.5"}
	cmd.Execute()
	os.Args = []string{"thicc", "modify", "0", "159.5"}
	cmd.Execute()
	os.Args = []string{"thicc", "modify", "1", "0"}
	cmd.Execute()

	// Test 'delete'
	os.Args = []string{"thicc", "delete", "2"}
	cmd.Execute()
	os.Args = []string{"thicc", "delete", "0"}
	cmd.Execute()

	// Test 'reset' - cancel
	stdinR, stdinW, _ = os.Pipe()
	os.Stdin = stdinR
	go func() {
		stdinW.Write([]byte("no\n"))
		stdinW.Close()
	}()
	os.Args = []string{"thicc", "reset"}
	cmd.Execute()

	// Test 'reset' - confirm
	stdout = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w

	// Simulate 'yes' confirmation
	stdinR, stdinW, _ = os.Pipe()
	os.Stdin = stdinR
	go func() {
		stdinW.Write([]byte("yes\n"))
		stdinW.Close()
	}()

	os.Args = []string{"thicc", "reset"}
	cmd.Execute()

	w.Close()
	os.Stdout = stdout
	os.Stdin = oldStdin
	buf.Reset()
	io.Copy(&buf, r)
	output = buf.String()
	if !strings.Contains(output, "All weight entries and settings have been deleted") {
		t.Error("Expected reset confirmation, got:", output)
	}
}
