package tests

import (
	"math"
	"testing"

	"github.com/tryonlinux/thicc/internal/calculator"
)

func TestCalculateBMI_MetricUnits(t *testing.T) {
	// Test BMI calculation with kg and cm
	weight := 70.0 // kg
	height := 175.0 // cm
	bmi := calculator.CalculateBMI(weight, height, "kg", "cm")

	// BMI = 70 / (1.75^2) = 22.86
	expected := 22.86
	if math.Abs(bmi-expected) > 0.01 {
		t.Errorf("Expected BMI %.2f, got %.2f", expected, bmi)
	}
}

func TestCalculateBMI_ImperialUnits(t *testing.T) {
	// Test BMI calculation with lbs and in
	weight := 154.0 // lbs (approximately 70 kg)
	height := 69.0  // inches (approximately 175 cm)
	bmi := calculator.CalculateBMI(weight, height, "lbs", "in")

	// BMI = (154 / 69^2) * 703 = 22.74
	expected := 22.74
	if math.Abs(bmi-expected) > 0.01 {
		t.Errorf("Expected BMI %.2f, got %.2f", expected, bmi)
	}
}

func TestCalculateBMI_MixedUnits_KgInches(t *testing.T) {
	// Test BMI calculation with kg and inches
	weight := 70.0 // kg
	height := 69.0 // inches
	bmi := calculator.CalculateBMI(weight, height, "kg", "in")

	// Convert inches to meters: 69 * 0.0254 = 1.7526 m
	// BMI = 70 / (1.7526^2) = 22.79
	expected := 22.79
	if math.Abs(bmi-expected) > 0.01 {
		t.Errorf("Expected BMI %.2f, got %.2f", expected, bmi)
	}
}

func TestCalculateBMI_MixedUnits_LbsCm(t *testing.T) {
	// Test BMI calculation with lbs and cm
	weight := 154.0 // lbs
	height := 175.0 // cm
	bmi := calculator.CalculateBMI(weight, height, "lbs", "cm")

	// Convert lbs to kg: 154 * 0.453592 = 69.85 kg
	// Convert cm to m: 175 / 100 = 1.75 m
	// BMI = 69.85 / (1.75^2) = 22.81
	expected := 22.81
	if math.Abs(bmi-expected) > 0.01 {
		t.Errorf("Expected BMI %.2f, got %.2f", expected, bmi)
	}
}

func TestCalculateBMI_Underweight(t *testing.T) {
	// Test underweight BMI
	weight := 50.0   // kg
	height := 175.0  // cm
	bmi := calculator.CalculateBMI(weight, height, "kg", "cm")

	// BMI = 50 / (1.75^2) = 16.33 (underweight)
	expected := 16.33
	if math.Abs(bmi-expected) > 0.01 {
		t.Errorf("Expected BMI %.2f, got %.2f", expected, bmi)
	}
}

func TestCalculateBMI_Overweight(t *testing.T) {
	// Test overweight BMI
	weight := 100.0  // kg
	height := 175.0  // cm
	bmi := calculator.CalculateBMI(weight, height, "kg", "cm")

	// BMI = 100 / (1.75^2) = 32.65 (overweight)
	expected := 32.65
	if math.Abs(bmi-expected) > 0.01 {
		t.Errorf("Expected BMI %.2f, got %.2f", expected, bmi)
	}
}
