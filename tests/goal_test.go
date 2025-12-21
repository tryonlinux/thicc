package tests

import (
	"testing"

	"github.com/tryonlinux/thicc/internal/models"
)

func TestGoalDifferenceCalculation(t *testing.T) {
	testCases := []struct {
		name          string
		currentWeight float64
		goalWeight    float64
		expectedDiff  float64
		shouldLose    bool
		shouldGain    bool
		atGoal        bool
	}{
		{
			name:          "Need to lose weight",
			currentWeight: 160.0,
			goalWeight:    150.0,
			expectedDiff:  10.0,
			shouldLose:    true,
		},
		{
			name:          "Need to gain weight",
			currentWeight: 140.0,
			goalWeight:    150.0,
			expectedDiff:  -10.0,
			shouldGain:    true,
		},
		{
			name:          "At goal weight",
			currentWeight: 150.0,
			goalWeight:    150.0,
			expectedDiff:  0.0,
			atGoal:        true,
		},
		{
			name:          "Small amount to lose",
			currentWeight: 151.5,
			goalWeight:    150.0,
			expectedDiff:  1.5,
			shouldLose:    true,
		},
		{
			name:          "Large amount to gain",
			currentWeight: 120.0,
			goalWeight:    160.0,
			expectedDiff:  -40.0,
			shouldGain:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			diff := tc.currentWeight - tc.goalWeight

			if diff != tc.expectedDiff {
				t.Errorf("Expected difference %.2f, got %.2f", tc.expectedDiff, diff)
			}

			if tc.shouldLose && diff <= 0 {
				t.Errorf("Expected to need to lose weight, but diff is %.2f", diff)
			}

			if tc.shouldGain && diff >= 0 {
				t.Errorf("Expected to need to gain weight, but diff is %.2f", diff)
			}

			if tc.atGoal && diff != 0 {
				t.Errorf("Expected to be at goal, but diff is %.2f", diff)
			}
		})
	}
}

func TestGoalWeightWithWeightEntries(t *testing.T) {
	db := setupTestDB(t)

	// Set up settings with goal weight
	db.Exec("INSERT INTO settings (key, value) VALUES ('weight_unit', 'lbs')")
	db.Exec("INSERT INTO settings (key, value) VALUES ('height_unit', 'in')")
	db.Exec("INSERT INTO settings (key, value) VALUES ('height', '70')")
	db.Exec("INSERT INTO settings (key, value) VALUES ('goal_weight', '150')")

	settings, err := models.GetSettings(db)
	if err != nil {
		t.Fatalf("Failed to get settings: %v", err)
	}

	// Add weight entries
	models.AddWeight(db, "2024-01-01", 160.0, 23.0)
	models.AddWeight(db, "2024-01-05", 158.0, 22.7)
	models.AddWeight(db, "2024-01-10", 155.0, 22.3)
	models.AddWeight(db, "2024-01-15", 152.0, 21.8)
	models.AddWeight(db, "2024-01-20", 150.0, 21.5) // At goal

	// Get latest weight
	weights, err := models.GetWeights(db, 1)
	if err != nil {
		t.Fatalf("Failed to get weights: %v", err)
	}

	latestWeight := weights[0].Weight

	// Calculate difference
	diff := latestWeight - settings.GoalWeight

	if diff != 0.0 {
		t.Errorf("Expected to be at goal (diff = 0), got diff = %.2f", diff)
	}

	// Add another entry above goal
	models.AddWeight(db, "2024-01-25", 155.0, 22.3)

	weights, err = models.GetWeights(db, 1)
	if err != nil {
		t.Fatalf("Failed to get weights: %v", err)
	}

	latestWeight = weights[0].Weight
	diff = latestWeight - settings.GoalWeight

	if diff != 5.0 {
		t.Errorf("Expected diff of 5.0 (need to lose), got %.2f", diff)
	}

	if diff <= 0 {
		t.Errorf("Should need to lose weight, but diff is %.2f", diff)
	}
}

func TestGoalWeightEdgeCases(t *testing.T) {
	db := setupTestDB(t)

	testCases := []struct {
		name       string
		goalWeight string
		valid      bool
	}{
		{"Normal goal", "150.5", true},
		{"Zero goal", "0", true}, // Technically valid but not realistic
		{"Large goal", "500.0", true},
		{"Small goal", "0.1", true},
		{"Decimal goal", "145.75", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear settings
			db.Exec("DELETE FROM settings")

			// Insert test settings
			db.Exec("INSERT INTO settings (key, value) VALUES ('weight_unit', 'lbs')")
			db.Exec("INSERT INTO settings (key, value) VALUES ('height_unit', 'in')")
			db.Exec("INSERT INTO settings (key, value) VALUES ('height', '70')")
			db.Exec("INSERT INTO settings (key, value) VALUES ('goal_weight', ?)", tc.goalWeight)

			settings, err := models.GetSettings(db)
			if tc.valid && err != nil {
				t.Errorf("Expected valid goal weight, got error: %v", err)
			}
			if tc.valid && settings == nil {
				t.Errorf("Expected settings to be returned")
			}
		})
	}
}
