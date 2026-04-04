package tests

import (
	"errors"
	"testing"

	"github.com/tryonlinux/thicc/internal/validation"
)

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name    string
		dateStr string
		wantErr error
	}{
		{"valid date", "2024-01-01", nil},
		{"empty date", "", validation.ErrInvalidDateFormat},
		{"invalid format", "01-01-2024", validation.ErrInvalidDate},
		{"invalid date", "2024-13-01", validation.ErrInvalidDate},
		{"invalid day", "2024-01-32", validation.ErrInvalidDate},
		{"not a date", "not-a-date", validation.ErrInvalidDate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateDate(tt.dateStr)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateDate(%q) error = %v, wantErr %v", tt.dateStr, err, tt.wantErr)
			}
		})
	}
}

func TestValidateWeight(t *testing.T) {
	tests := []struct {
		name    string
		weight  float64
		wantErr error
	}{
		{"valid weight", 70.0, nil},
		{"min weight", 1.0, nil},
		{"max weight", 1000.0, nil},
		{"zero weight", 0.0, validation.ErrNegativeNumber},
		{"negative weight", -1.0, validation.ErrNegativeNumber},
		{"too low weight", 0.5, validation.ErrInvalidWeight},
		{"too high weight", 1001.0, validation.ErrInvalidWeight},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateWeight(tt.weight)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateWeight(%v) error = %v, wantErr %v", tt.weight, err, tt.wantErr)
			}
		})
	}
}

func TestValidateBMI(t *testing.T) {
	tests := []struct {
		name    string
		bmi     float64
		wantErr error
	}{
		{"valid BMI", 22.5, nil},
		{"min BMI", 5.0, nil},
		{"max BMI", 100.0, nil},
		{"too low BMI", 4.9, validation.ErrInvalidBMI},
		{"too high BMI", 100.1, validation.ErrInvalidBMI},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateBMI(tt.bmi)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateBMI(%v) error = %v, wantErr %v", tt.bmi, err, tt.wantErr)
			}
		})
	}
}

func TestValidateHeight(t *testing.T) {
	tests := []struct {
		name    string
		height  float64
		unit    string
		wantErr error
	}{
		{"valid cm", 175.0, "cm", nil},
		{"min cm", 50.0, "cm", nil},
		{"max cm", 300.0, "cm", nil},
		{"too low cm", 49.9, "cm", validation.ErrInvalidHeightCm},
		{"too high cm", 300.1, "cm", validation.ErrInvalidHeightCm},
		{"valid in", 70.0, "in", nil},
		{"min in", 20.0, "in", nil},
		{"max in", 120.0, "in", nil},
		{"too low in", 19.9, "in", validation.ErrInvalidHeightIn},
		{"too high in", 120.1, "in", validation.ErrInvalidHeightIn},
		{"zero height", 0.0, "cm", validation.ErrNegativeNumber},
		{"negative height", -1.0, "in", validation.ErrNegativeNumber},
		{"unknown unit", 175.0, "m", nil}, // Based on current implementation, it just returns nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateHeight(tt.height, tt.unit)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateHeight(%v, %q) error = %v, wantErr %v", tt.height, tt.unit, err, tt.wantErr)
			}
		})
	}
}

func TestParsePositiveFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"valid float", "70.5", 70.5, false},
		{"valid float with spaces", "  154.3  ", 154.3, false},
		{"empty string", "", 0, true},
		{"just spaces", "   ", 0, true},
		{"not a number", "abc", 0, true},
		{"zero", "0", 0, true},
		{"negative", "-5", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validation.ParsePositiveFloat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePositiveFloat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePositiveFloat(%q) got = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseAndValidateWeight(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"valid weight", "70", 70, false},
		{"invalid format", "abc", 0, true},
		{"out of range low", "0.5", 0, true},
		{"out of range high", "2000", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validation.ParseAndValidateWeight(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAndValidateWeight(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseAndValidateWeight(%q) got = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
