package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tryonlinux/thicc/internal/validation"
)

var goalCmd = &cobra.Command{
	Use:   "goal <weight>",
	Short: "Set your goal weight",
	Long:  `Set or update your goal weight target.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()
		settings := GetSettings()

		// Parse and validate goal weight
		goalWeight, err := validation.ParseAndValidateWeight(args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Update in database
		_, err = db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES ('goal_weight', ?)",
			strconv.FormatFloat(goalWeight, 'f', 2, 64))
		if err != nil {
			fmt.Printf("Error updating goal weight: %v\n", err)
			return
		}

		// Update settings in memory
		settings.GoalWeight = goalWeight

		fmt.Printf("Goal weight set to %.2f %s\n", goalWeight, settings.WeightUnit)

		// Show updated table
		showCmd.Run(cmd, []string{})
	},
}
