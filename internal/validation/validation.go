package validation

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Date format constants
const DateFormat = "2006-01-02"

// Weight bounds (in any unit)
const (
	MinWeight = 1.0
	MaxWeight = 1000.0
)

// BMI bounds
const (
	MinBMI = 5.0
	MaxBMI = 100.0
)

// Height bounds
const (
	MinHeightCm = 50.0
	MaxHeightCm = 300.0
	MinHeightIn = 20.0
	MaxHeightIn = 120.0
)

// Common errors
var (
	ErrInvalidWeight     = errors.New("weight must be between 1 and 1000")
	ErrInvalidBMI        = errors.New("BMI must be between 5 and 100")
	ErrInvalidHeightCm   = errors.New("height must be between 50 and 300 cm")
	ErrInvalidHeightIn   = errors.New("height must be between 20 and 120 inches")
	ErrInvalidDate       = errors.New("date must be in YYYY-MM-DD format and be a valid date")
	ErrNegativeNumber    = errors.New("value must be a positive number")
	ErrInvalidDateFormat = errors.New("date format must be YYYY-MM-DD")
)

// ValidateDate validates a date string is in YYYY-MM-DD format and is a valid date
func ValidateDate(dateStr string) error {
	if dateStr == "" {
		return ErrInvalidDateFormat
	}

	// Parse the date using time.Parse to ensure it's a valid date
	_, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return ErrInvalidDate
	}

	return nil
}

// ValidateWeight validates a weight value is within reasonable bounds
func ValidateWeight(weight float64) error {
	if weight <= 0 {
		return ErrNegativeNumber
	}
	if weight < MinWeight || weight > MaxWeight {
		return ErrInvalidWeight
	}
	return nil
}

// ValidateBMI validates a BMI value is within reasonable bounds
func ValidateBMI(bmi float64) error {
	if bmi < MinBMI || bmi > MaxBMI {
		return ErrInvalidBMI
	}
	return nil
}

// ValidateHeight validates height based on unit
func ValidateHeight(height float64, unit string) error {
	if height <= 0 {
		return ErrNegativeNumber
	}

	if unit == "cm" {
		if height < MinHeightCm || height > MaxHeightCm {
			return ErrInvalidHeightCm
		}
	} else if unit == "in" {
		if height < MinHeightIn || height > MaxHeightIn {
			return ErrInvalidHeightIn
		}
	}

	return nil
}

// ParsePositiveFloat parses a string to float64, trims whitespace, and validates it's positive
func ParsePositiveFloat(s string) (float64, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return 0, ErrNegativeNumber
	}

	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, errors.New("invalid number format")
	}

	if value <= 0 {
		return 0, ErrNegativeNumber
	}

	return value, nil
}

// ParseAndValidateWeight parses and validates a weight value in one step
func ParseAndValidateWeight(s string) (float64, error) {
	weight, err := ParsePositiveFloat(s)
	if err != nil {
		return 0, err
	}

	if err := ValidateWeight(weight); err != nil {
		return 0, err
	}

	return weight, nil
}
