package models

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tryonlinux/thicc/internal/database"
)

// Settings represents application settings
type Settings struct {
	WeightUnit string  // "lbs" or "kg"
	HeightUnit string  // "in" or "cm"
	Height     float64 // height in the specified unit
	GoalWeight float64 // goal weight in the specified unit
}

// GetSettings retrieves current application settings
func GetSettings(db *database.DB) (*Settings, error) {
	var weightUnit, heightUnit, heightStr, goalWeightStr string

	err := db.QueryRow("SELECT value FROM settings WHERE key = 'weight_unit'").Scan(&weightUnit)
	if err == sql.ErrNoRows {
		// First launch - need to setup
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	err = db.QueryRow("SELECT value FROM settings WHERE key = 'height_unit'").Scan(&heightUnit)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow("SELECT value FROM settings WHERE key = 'height'").Scan(&heightStr)
	if err != nil {
		return nil, err
	}

	height, err := strconv.ParseFloat(heightStr, 64)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow("SELECT value FROM settings WHERE key = 'goal_weight'").Scan(&goalWeightStr)
	if err != nil {
		return nil, err
	}

	goalWeight, err := strconv.ParseFloat(goalWeightStr, 64)
	if err != nil {
		return nil, err
	}

	return &Settings{
		WeightUnit: weightUnit,
		HeightUnit: heightUnit,
		Height:     height,
		GoalWeight: goalWeight,
	}, nil
}

// SetupSettings prompts the user for initial settings
func SetupSettings(db *database.DB) (*Settings, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== First Time Setup ===")
	fmt.Println("Please configure your preferences.\n")

	// Get weight unit
	var weightUnit string
	for {
		fmt.Print("Weight unit (lbs/kg): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		weightUnit = strings.TrimSpace(strings.ToLower(input))
		if weightUnit == "lbs" || weightUnit == "kg" {
			break
		}
		fmt.Println("Invalid input. Please enter 'lbs' or 'kg'.")
	}

	// Get height unit
	var heightUnit string
	for {
		fmt.Print("Height unit (in/cm): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		heightUnit = strings.TrimSpace(strings.ToLower(input))
		if heightUnit == "in" || heightUnit == "cm" {
			break
		}
		fmt.Println("Invalid input. Please enter 'in' or 'cm'.")
	}

	// Get height
	var height float64
	for {
		fmt.Printf("Your height (%s): ", heightUnit)
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		input = strings.TrimSpace(input)
		height, err = strconv.ParseFloat(input, 64)
		if err == nil && height > 0 {
			break
		}
		fmt.Println("Invalid input. Please enter a positive number.")
	}

	// Get goal weight
	var goalWeight float64
	for {
		fmt.Printf("Your goal weight (%s): ", weightUnit)
		input, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		input = strings.TrimSpace(input)
		goalWeight, err = strconv.ParseFloat(input, 64)
		if err == nil && goalWeight > 0 {
			break
		}
		fmt.Println("Invalid input. Please enter a positive number.")
	}

	// Save to database
	_, err := db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('weight_unit', ?)", weightUnit)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('height_unit', ?)", heightUnit)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('height', ?)", strconv.FormatFloat(height, 'f', 2, 64))
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('goal_weight', ?)", strconv.FormatFloat(goalWeight, 'f', 2, 64))
	if err != nil {
		return nil, err
	}

	fmt.Println("\nSettings saved successfully!\n")

	return &Settings{
		WeightUnit: weightUnit,
		HeightUnit: heightUnit,
		Height:     height,
		GoalWeight: goalWeight,
	}, nil
}

// ResetSettings clears all settings (used by reset command)
func ResetSettings(db *database.DB) error {
	_, err := db.Exec("DELETE FROM settings")
	return err
}
