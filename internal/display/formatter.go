package display

import "fmt"

// FormatWeight formats a weight value with proper precision and unit
func FormatWeight(weight float64, unit string) string {
	return fmt.Sprintf("%.2f %s", weight, unit)
}

// FormatBMI formats a BMI value with proper precision
func FormatBMI(bmi float64) string {
	return fmt.Sprintf("%.1f", bmi)
}

// FormatDate returns the date as-is (already in YYYY-MM-DD format)
func FormatDate(date string) string {
	return date
}
