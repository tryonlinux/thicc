# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

THICC is a weight tracking CLI tool built with Go that visualizes weight progress over time with tables and ASCII line graphs. Uses SQLite for local storage and Cobra for CLI framework.

## Build and Development Commands

```bash
# Build executable
go build -o thicc.exe

# Run tests
go test ./tests/... -v

# Run specific test
go test ./tests/... -v -run TestGoalWeightSetting

# Download dependencies
go mod tidy

# Seed demo data (requires built executable)
./setup-demo.sh   # Unix/Linux/macOS
setup-demo.bat    # Windows
```

## Architecture

### Core Flow
1. **Initialization** (`cmd/root.go`): On startup, `initDatabase()` runs via Cobra's `OnInitialize`
   - Opens database at `~/.thicc/weights.db`
   - Checks if settings exist (`models.GetSettings()`)
   - If no settings (first launch), runs `models.SetupSettings()` which prompts for units, height, and goal weight
   - Settings are stored in package-level variables `db` and `settings` accessible via `GetDB()` and `GetSettings()`

2. **Command Pattern**: All commands in `cmd/` access shared DB and settings via `GetDB()` and `GetSettings()`
   - Commands call model functions to perform database operations
   - After mutation commands (add, modify, delete, goal), `showCmd.Run()` is called to display updated data
   - Running just `thicc` defaults to `show` command

3. **Display Pipeline** (`internal/display/table.go`):
   - `RenderWeightsTable()` orchestrates the entire output
   - Calculates statistics (min, max, avg) from weight slice
   - Computes goal difference: `latestWeight - goalWeight`
     - Positive = need to lose
     - Negative = need to gain
   - Creates table (left) and graph (right) separately, then joins horizontally with lipgloss
   - Graph uses Bresenham's line algorithm to connect data points

### Data Layer

**Settings** (`internal/models/settings.go`):
- Stored as key-value pairs in `settings` table
- Keys: `weight_unit`, `height_unit`, `height`, `goal_weight`
- `GetSettings()` returns nil on first launch (triggers setup flow)

**Weights** (`internal/models/weight.go`):
- Each entry has: ID (auto-increment), date (YYYY-MM-DD), weight, BMI
- BMI is **calculated once on add** and stored (not recalculated on retrieval)
- Queries always return descending by date: `ORDER BY date DESC, id DESC`

**BMI Calculation** (`internal/calculator/bmi.go`):
- Supports 4 unit combinations: kg+cm, lbs+in, kg+in, lbs+cm
- Called by `add` and `modify` commands to compute BMI before storage

### Graph Rendering Details

The ASCII line graph (`createLineGraph()` in `internal/display/table.go`):
- Normalizes weights to fit 40x20 character grid
- Includes goal weight in min/max range calculation to ensure goal line is visible
- Characters used:
  - Data points: `·` (middle dot - smallest)
  - Connecting lines: `∙` (bullet operator - lighter)
  - Goal line: `─` (horizontal box drawing)
- Goal line drawn at calculated Y position with label "Goal: X.X" on left axis
- Reverses weight order to show oldest→newest (left→right)

### Database Schema

```sql
weights (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  date TEXT NOT NULL,
  weight REAL NOT NULL,
  bmi REAL NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
)

settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

## Key Design Decisions

1. **BMI Storage**: BMI is calculated once and stored (not computed on-the-fly) because it depends on user's height which is a setting, not per-weight data
2. **First Launch Detection**: Uses `sql.ErrNoRows` when querying settings to detect first launch
3. **Date Format**: All dates are YYYY-MM-DD strings (Go format: "2006-01-02")
4. **Graph Range**: Goal weight is always included in min/max calculation to ensure the goal line appears on the graph
5. **Table Truncation**: Show command can fetch unlimited weights for graphing, but table display always truncates to 20 entries
6. **Default Command**: Running `thicc` with no args shows the table (mimics `thicc show`)

## Testing

Tests are in `tests/` directory:
- `calculator_test.go`: BMI calculations across all unit combinations
- `models_test.go`: Database CRUD operations, settings management, goal weight functionality
- `goal_test.go`: Goal difference calculations and integration tests
- `display_test.go`: Formatters for weight, BMI, dates

All tests use `setupTestDB()` helper which creates temporary SQLite file that auto-cleans on test completion.

## GitHub Actions

- **PR Tests** (`.github/workflows/pr-tests.yml`): Runs tests and comments results on PRs
- **Release** (`.github/workflows/release.yml`):
  - Semantic versioning based on commit messages (BREAKING/feat/patch)
  - Cross-compiles for Windows/macOS/Linux (amd64, arm64, 386)
  - Creates GitHub release with all binaries
  - Runs on push to main

## Demo Data

`setup-demo.sh` and `setup-demo.bat` scripts:
- Reset database
- Configure with imperial units (lbs/in, 70" height, 145 lbs goal)
- Seed ~27 weight entries spanning Jan-May showing realistic weight loss journey
- Used for screenshots and testing visualizations
