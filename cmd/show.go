package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tryonlinux/thicc/internal/display"
	"github.com/tryonlinux/thicc/internal/models"
	"github.com/tryonlinux/thicc/internal/validation"
)

var showCmd = &cobra.Command{
	Use:   "show [number|date]",
	Short: "Display weight table and graph",
	Long: `Shows weight entries with a table and line graph.

Examples:
  thicc show          # Show last 20 entries
  thicc show 50       # Show last 50 entries
  thicc show 2024-01-01  # Show entries from 2024-01-01 to today (table shows last 20)`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()
		settings := GetSettings()

		var weights []models.Weight
		var err error
		limit := display.DefaultDisplayLimit

		if len(args) == 0 {
			// Default: show last entries
			weights, err = models.GetWeights(db, display.DefaultDisplayLimit)
		} else {
			arg := strings.TrimSpace(args[0])

			// Check if it's a number (limit) or date
			if num, err := strconv.Atoi(arg); err == nil {
				// It's a number
				limit = num
				if limit <= 0 {
					fmt.Println("Error: Number must be positive")
					return
				}
				weights, err = models.GetWeights(db, limit)
				if err != nil {
					fmt.Printf("Error retrieving weights: %v\n", err)
					return
				}
			} else {
				// Try to parse as a date
				if err := validation.ValidateDate(arg); err == nil {
					// It's a valid date
					startDate := arg
					endDate := models.GetTodayDate()
					weights, err = models.GetWeightsBetweenDates(db, startDate, endDate)
					if err != nil {
						fmt.Printf("Error retrieving weights: %v\n", err)
						return
					}
					// For graph, use all weights; for table display, it will be truncated in render
				} else {
					fmt.Println("Error: Argument must be a positive number or a date in YYYY-MM-DD format")
					return
				}
			}
		}

		if err != nil {
			fmt.Printf("Error retrieving weights: %v\n", err)
			return
		}

		// Render table and graph
		output := display.RenderWeightsTable(weights, settings, limit)
		fmt.Println(output)
	},
}
