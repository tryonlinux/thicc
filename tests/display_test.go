package tests

import (
	"strings"
	"testing"

	"github.com/tryonlinux/thicc/internal/display"
	"github.com/tryonlinux/thicc/internal/models"
)

func TestFormatWeight(t *testing.T) {
	tests := []struct {
		weight   float64
		unit     string
		expected string
	}{
		{70.5, "kg", "70.50 kg"},
		{154.32, "lbs", "154.32 lbs"},
		{100.0, "kg", "100.00 kg"},
	}

	for _, tt := range tests {
		result := display.FormatWeight(tt.weight, tt.unit)
		if result != tt.expected {
			t.Errorf("FormatWeight(%.2f, %s) = %s, want %s", tt.weight, tt.unit, result, tt.expected)
		}
	}
}

func TestFormatBMI(t *testing.T) {
	tests := []struct {
		bmi      float64
		expected string
	}{
		{22.86, "22.9"},
		{18.5, "18.5"},
		{30.12, "30.1"},
		{25.0, "25.0"},
	}

	for _, tt := range tests {
		result := display.FormatBMI(tt.bmi)
		if result != tt.expected {
			t.Errorf("FormatBMI(%.2f) = %s, want %s", tt.bmi, result, tt.expected)
		}
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		date     string
		expected string
	}{
		{"2024-01-01", "2024-01-01"},
		{"2023-12-25", "2023-12-25"},
	}

	for _, tt := range tests {
		result := display.FormatDate(tt.date)
		if result != tt.expected {
			t.Errorf("FormatDate(%s) = %s, want %s", tt.date, result, tt.expected)
		}
	}
}

func TestRenderWeightsTableDelta(t *testing.T) {
	tests := []struct {
		name        string
		weights     []models.Weight
		settings    *models.Settings
		expectedStr string
	}{
		{
			name: "weight loss shows 'Lost'",
			weights: []models.Weight{
				{ID: 3, Date: "2024-01-15", Weight: 150.0, BMI: 22.0},
				{ID: 2, Date: "2024-01-08", Weight: 155.0, BMI: 22.5},
				{ID: 1, Date: "2024-01-01", Weight: 160.0, BMI: 23.0},
			},
			settings: &models.Settings{
				WeightUnit: "lbs",
				HeightUnit: "in",
				Height:     70,
				GoalWeight: 145,
			},
			expectedStr: "Lost 10.00 lbs",
		},
		{
			name: "weight gain shows 'Gained'",
			weights: []models.Weight{
				{ID: 3, Date: "2024-01-15", Weight: 165.0, BMI: 23.5},
				{ID: 2, Date: "2024-01-08", Weight: 160.0, BMI: 23.0},
				{ID: 1, Date: "2024-01-01", Weight: 155.0, BMI: 22.5},
			},
			settings: &models.Settings{
				WeightUnit: "lbs",
				HeightUnit: "in",
				Height:     70,
				GoalWeight: 150,
			},
			expectedStr: "Gained 10.00 lbs",
		},
		{
			name: "no change shows 'No change'",
			weights: []models.Weight{
				{ID: 2, Date: "2024-01-08", Weight: 160.0, BMI: 23.0},
				{ID: 1, Date: "2024-01-01", Weight: 160.0, BMI: 23.0},
			},
			settings: &models.Settings{
				WeightUnit: "lbs",
				HeightUnit: "in",
				Height:     70,
				GoalWeight: 150,
			},
			expectedStr: "No change",
		},
		{
			name: "metric units show correctly",
			weights: []models.Weight{
				{ID: 2, Date: "2024-01-08", Weight: 70.0, BMI: 22.0},
				{ID: 1, Date: "2024-01-01", Weight: 75.0, BMI: 23.0},
			},
			settings: &models.Settings{
				WeightUnit: "kg",
				HeightUnit: "cm",
				Height:     180,
				GoalWeight: 68,
			},
			expectedStr: "Lost 5.00 kg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := display.RenderWeightsTable(tt.weights, tt.settings, 20)
			if !strings.Contains(result, tt.expectedStr) {
				t.Errorf("RenderWeightsTable() output does not contain expected delta string '%s'", tt.expectedStr)
			}
		})
	}
}
