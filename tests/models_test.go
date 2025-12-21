package tests

import (
	"os"
	"testing"

	"github.com/tryonlinux/thicc/internal/database"
	"github.com/tryonlinux/thicc/internal/models"
)

// setupTestDB creates a temporary database for testing
func setupTestDB(t *testing.T) *database.DB {
	// Create a temporary database file
	tmpFile, err := os.CreateTemp("", "thicc_test_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp database: %v", err)
	}
	tmpFile.Close()

	db, err := database.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Clean up function
	t.Cleanup(func() {
		db.Close()
		os.Remove(tmpFile.Name())
	})

	return db
}

func TestAddAndGetWeights(t *testing.T) {
	db := setupTestDB(t)

	// Add some test weights
	err := models.AddWeight(db, "2024-01-01", 70.0, 22.8)
	if err != nil {
		t.Fatalf("Failed to add weight: %v", err)
	}

	err = models.AddWeight(db, "2024-01-02", 69.5, 22.6)
	if err != nil {
		t.Fatalf("Failed to add weight: %v", err)
	}

	// Get weights
	weights, err := models.GetWeights(db, 10)
	if err != nil {
		t.Fatalf("Failed to get weights: %v", err)
	}

	if len(weights) != 2 {
		t.Errorf("Expected 2 weights, got %d", len(weights))
	}

	// Check they're in descending order by date
	if weights[0].Date != "2024-01-02" {
		t.Errorf("Expected first weight to be from 2024-01-02, got %s", weights[0].Date)
	}
}

func TestDeleteWeight(t *testing.T) {
	db := setupTestDB(t)

	// Add a weight
	err := models.AddWeight(db, "2024-01-01", 70.0, 22.8)
	if err != nil {
		t.Fatalf("Failed to add weight: %v", err)
	}

	// Get the weight ID
	weights, err := models.GetWeights(db, 1)
	if err != nil || len(weights) == 0 {
		t.Fatalf("Failed to get weight")
	}

	id := weights[0].ID

	// Delete the weight
	err = models.DeleteWeight(db, id)
	if err != nil {
		t.Fatalf("Failed to delete weight: %v", err)
	}

	// Verify it's gone
	weights, err = models.GetWeights(db, 10)
	if err != nil {
		t.Fatalf("Failed to get weights: %v", err)
	}

	if len(weights) != 0 {
		t.Errorf("Expected 0 weights after deletion, got %d", len(weights))
	}
}

func TestModifyWeight(t *testing.T) {
	db := setupTestDB(t)

	// Add a weight
	err := models.AddWeight(db, "2024-01-01", 70.0, 22.8)
	if err != nil {
		t.Fatalf("Failed to add weight: %v", err)
	}

	// Get the weight ID
	weights, err := models.GetWeights(db, 1)
	if err != nil || len(weights) == 0 {
		t.Fatalf("Failed to get weight")
	}

	id := weights[0].ID

	// Modify the weight
	newWeight := 65.0
	newBMI := 21.2
	err = models.ModifyWeight(db, id, newWeight, newBMI)
	if err != nil {
		t.Fatalf("Failed to modify weight: %v", err)
	}

	// Verify the change
	weights, err = models.GetWeights(db, 1)
	if err != nil || len(weights) == 0 {
		t.Fatalf("Failed to get weight after modification")
	}

	if weights[0].Weight != newWeight {
		t.Errorf("Expected weight %.2f, got %.2f", newWeight, weights[0].Weight)
	}

	if weights[0].BMI != newBMI {
		t.Errorf("Expected BMI %.2f, got %.2f", newBMI, weights[0].BMI)
	}
}

func TestGetWeightsBetweenDates(t *testing.T) {
	db := setupTestDB(t)

	// Add weights on different dates
	models.AddWeight(db, "2024-01-01", 70.0, 22.8)
	models.AddWeight(db, "2024-01-05", 69.5, 22.6)
	models.AddWeight(db, "2024-01-10", 69.0, 22.4)
	models.AddWeight(db, "2024-01-15", 68.5, 22.2)

	// Get weights between Jan 5 and Jan 12
	weights, err := models.GetWeightsBetweenDates(db, "2024-01-05", "2024-01-12")
	if err != nil {
		t.Fatalf("Failed to get weights between dates: %v", err)
	}

	if len(weights) != 2 {
		t.Errorf("Expected 2 weights, got %d", len(weights))
	}

	// Verify the dates
	if weights[0].Date != "2024-01-10" && weights[1].Date != "2024-01-10" {
		t.Errorf("Expected to find weight from 2024-01-10")
	}

	if weights[0].Date != "2024-01-05" && weights[1].Date != "2024-01-05" {
		t.Errorf("Expected to find weight from 2024-01-05")
	}
}

func TestGoalWeightSetting(t *testing.T) {
	db := setupTestDB(t)

	// Set up initial settings
	_, err := db.Exec("INSERT INTO settings (key, value) VALUES ('weight_unit', 'lbs')")
	if err != nil {
		t.Fatalf("Failed to insert weight_unit: %v", err)
	}
	_, err = db.Exec("INSERT INTO settings (key, value) VALUES ('height_unit', 'in')")
	if err != nil {
		t.Fatalf("Failed to insert height_unit: %v", err)
	}
	_, err = db.Exec("INSERT INTO settings (key, value) VALUES ('height', '70')")
	if err != nil {
		t.Fatalf("Failed to insert height: %v", err)
	}
	_, err = db.Exec("INSERT INTO settings (key, value) VALUES ('goal_weight', '150')")
	if err != nil {
		t.Fatalf("Failed to insert goal_weight: %v", err)
	}

	// Get settings
	settings, err := models.GetSettings(db)
	if err != nil {
		t.Fatalf("Failed to get settings: %v", err)
	}

	if settings.GoalWeight != 150.0 {
		t.Errorf("Expected goal weight 150.0, got %.2f", settings.GoalWeight)
	}

	if settings.WeightUnit != "lbs" {
		t.Errorf("Expected weight unit 'lbs', got '%s'", settings.WeightUnit)
	}

	if settings.HeightUnit != "in" {
		t.Errorf("Expected height unit 'in', got '%s'", settings.HeightUnit)
	}

	if settings.Height != 70.0 {
		t.Errorf("Expected height 70.0, got %.2f", settings.Height)
	}
}

func TestUpdateGoalWeight(t *testing.T) {
	db := setupTestDB(t)

	// Set up initial settings
	_, err := db.Exec("INSERT INTO settings (key, value) VALUES ('weight_unit', 'kg')")
	if err != nil {
		t.Fatalf("Failed to insert weight_unit: %v", err)
	}
	_, err = db.Exec("INSERT INTO settings (key, value) VALUES ('height_unit', 'cm')")
	if err != nil {
		t.Fatalf("Failed to insert height_unit: %v", err)
	}
	_, err = db.Exec("INSERT INTO settings (key, value) VALUES ('height', '175')")
	if err != nil {
		t.Fatalf("Failed to insert height: %v", err)
	}
	_, err = db.Exec("INSERT INTO settings (key, value) VALUES ('goal_weight', '70')")
	if err != nil {
		t.Fatalf("Failed to insert goal_weight: %v", err)
	}

	// Update goal weight
	newGoal := 65.0
	_, err = db.Exec("UPDATE settings SET value = ? WHERE key = 'goal_weight'", "65")
	if err != nil {
		t.Fatalf("Failed to update goal weight: %v", err)
	}

	// Get settings again
	settings, err := models.GetSettings(db)
	if err != nil {
		t.Fatalf("Failed to get settings: %v", err)
	}

	if settings.GoalWeight != newGoal {
		t.Errorf("Expected updated goal weight %.2f, got %.2f", newGoal, settings.GoalWeight)
	}
}

func TestGoalWeightDifferentUnits(t *testing.T) {
	db := setupTestDB(t)

	testCases := []struct {
		name       string
		weightUnit string
		heightUnit string
		height     string
		goalWeight string
		expected   float64
	}{
		{"Imperial", "lbs", "in", "70", "150", 150.0},
		{"Metric", "kg", "cm", "175", "70", 70.0},
		{"Mixed LbsCm", "lbs", "cm", "178", "165", 165.0},
		{"Mixed KgIn", "kg", "in", "69", "68", 68.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear settings
			db.Exec("DELETE FROM settings")

			// Insert test settings
			db.Exec("INSERT INTO settings (key, value) VALUES ('weight_unit', ?)", tc.weightUnit)
			db.Exec("INSERT INTO settings (key, value) VALUES ('height_unit', ?)", tc.heightUnit)
			db.Exec("INSERT INTO settings (key, value) VALUES ('height', ?)", tc.height)
			db.Exec("INSERT INTO settings (key, value) VALUES ('goal_weight', ?)", tc.goalWeight)

			settings, err := models.GetSettings(db)
			if err != nil {
				t.Fatalf("Failed to get settings: %v", err)
			}

			if settings.GoalWeight != tc.expected {
				t.Errorf("Expected goal weight %.2f, got %.2f", tc.expected, settings.GoalWeight)
			}

			if settings.WeightUnit != tc.weightUnit {
				t.Errorf("Expected weight unit '%s', got '%s'", tc.weightUnit, settings.WeightUnit)
			}
		})
	}
}
