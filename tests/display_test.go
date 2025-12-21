package tests

import (
	"testing"

	"github.com/tryonlinux/thicc/internal/display"
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
