package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tryonlinux/thicc/internal/models"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Clear all data and start over",
	Long:  `Deletes all weight entries and settings. You will be prompted to reconfigure on next launch. This action cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()

		// Get confirmation from user
		fmt.Println("WARNING: This will delete ALL weight entries and settings.")
		fmt.Print("Are you sure you want to continue? (yes/no): ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "yes" && response != "y" {
			fmt.Println("Reset cancelled.")
			return
		}

		// Delete all weights
		_, err = db.Exec("DELETE FROM weights")
		if err != nil {
			fmt.Printf("Error deleting weights: %v\n", err)
			return
		}

		// Reset settings
		err = models.ResetSettings(db)
		if err != nil {
			fmt.Printf("Error resetting settings: %v\n", err)
			return
		}

		fmt.Println("\nAll weight entries and settings have been deleted.")
		fmt.Println("You will be prompted to reconfigure on next launch.")
	},
}
