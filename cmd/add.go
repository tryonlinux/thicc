package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tryonlinux/thicc/internal/calculator"
	"github.com/tryonlinux/thicc/internal/models"
	"github.com/tryonlinux/thicc/internal/validation"
)

var addCmd = &cobra.Command{
	Use:   "add <weight> [date]",
	Short: "Add a new weight entry",
	Long:  `Add a new weight entry with optional date (defaults to today). Date format: YYYY-MM-DD`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()
		settings := GetSettings()

		// Parse and validate weight
		weight, err := validation.ParseAndValidateWeight(args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Parse date (default to today)
		date := models.GetTodayDate()
		if len(args) == 2 {
			date = strings.TrimSpace(args[1])
			// Validate date format and validity
			if err := validation.ValidateDate(date); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}

		// Calculate BMI
		bmi := calculator.CalculateBMI(weight, settings.Height, settings.WeightUnit, settings.HeightUnit)

		// Add to database
		err = models.AddWeight(db, date, weight, bmi)
		if err != nil {
			fmt.Printf("Error adding weight: %v\n", err)
			return
		}

		fmt.Printf("Added weight: %.2f %s on %s (BMI: %.1f)\n", weight, settings.WeightUnit, date, bmi)

		// Show updated table
		showCmd.Run(cmd, []string{})
	},
}
