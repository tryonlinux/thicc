package calculator

// CalculateBMI calculates BMI based on weight, height, and units
func CalculateBMI(weight, height float64, weightUnit, heightUnit string) float64 {
	var bmi float64

	if weightUnit == "kg" && heightUnit == "cm" {
		// BMI = kg / (m^2)
		heightInMeters := height / 100.0
		bmi = weight / (heightInMeters * heightInMeters)
	} else if weightUnit == "lbs" && heightUnit == "in" {
		// BMI = (lbs / in^2) * 703
		bmi = (weight / (height * height)) * 703
	} else if weightUnit == "kg" && heightUnit == "in" {
		// Convert inches to meters
		heightInMeters := height * 0.0254
		bmi = weight / (heightInMeters * heightInMeters)
	} else { // lbs and cm
		// Convert lbs to kg and cm to meters
		weightInKg := weight * 0.453592
		heightInMeters := height / 100.0
		bmi = weightInKg / (heightInMeters * heightInMeters)
	}

	return bmi
}
