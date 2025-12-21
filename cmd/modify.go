package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tryonlinux/thicc/internal/calculator"
	"github.com/tryonlinux/thicc/internal/models"
	"github.com/tryonlinux/thicc/internal/validation"
)

var modifyCmd = &cobra.Command{
	Use:   "modify <weightId> <weight>",
	Short: "Modify a weight entry",
	Long:  `Modify a weight entry by its ID (shown in the show command).`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()
		settings := GetSettings()

		// Parse weight ID
		id, err := strconv.Atoi(strings.TrimSpace(args[0]))
		if err != nil || id <= 0 {
			fmt.Println("Error: Weight ID must be a positive number")
			return
		}

		// Parse and validate weight
		weight, err := validation.ParseAndValidateWeight(args[1])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Calculate new BMI
		bmi := calculator.CalculateBMI(weight, settings.Height, settings.WeightUnit, settings.HeightUnit)

		// Update in database
		err = models.ModifyWeight(db, id, weight, bmi)
		if err != nil {
			fmt.Printf("Error modifying weight: %v\n", err)
			return
		}

		fmt.Printf("Updated weight entry %d to %.2f %s (BMI: %.1f)\n", id, weight, settings.WeightUnit, bmi)

		// Show updated table
		showCmd.Run(cmd, []string{})
	},
}
