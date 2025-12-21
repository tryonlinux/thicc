package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tryonlinux/thicc/internal/models"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <weightId>",
	Short: "Delete a weight entry",
	Long:  `Delete a weight entry by its ID (shown in the show command).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()

		// Parse weight ID
		id, err := strconv.Atoi(strings.TrimSpace(args[0]))
		if err != nil || id <= 0 {
			fmt.Println("Error: Weight ID must be a positive number")
			return
		}

		// Delete from database
		err = models.DeleteWeight(db, id)
		if err != nil {
			fmt.Printf("Error deleting weight: %v\n", err)
			return
		}

		fmt.Printf("Deleted weight entry with ID %d\n", id)

		// Show updated table
		showCmd.Run(cmd, []string{})
	},
}
