package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tryonlinux/thicc/internal/config"
	"github.com/tryonlinux/thicc/internal/database"
	"github.com/tryonlinux/thicc/internal/models"
)

var db *database.DB
var settings *models.Settings

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "thicc",
	Short: "THICC - Weight tracking CLI",
	Long:  `THICC helps you track your weight and visualize your progress over time.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior: run show command
		if settings != nil {
			showCmd.Run(cmd, args)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Clean up database connection on exit
		cleanupDatabase()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initDatabase)

	// Add all subcommands
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(modifyCmd)
	rootCmd.AddCommand(goalCmd)
	rootCmd.AddCommand(resetCmd)
}

// initDatabase initializes the database connection
func initDatabase() {
	dbPath, err := config.GetDatabasePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting database path: %v\n", err)
		os.Exit(1)
	}

	db, err = database.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}

	// Check if settings exist (first launch detection)
	settings, err = models.GetSettings(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting settings: %v\n", err)
		os.Exit(1)
	}

	// First launch - prompt for setup
	if settings == nil {
		settings, err = models.SetupSettings(db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error setting up: %v\n", err)
			os.Exit(1)
		}
	}
}

// GetDB returns the database connection (used by commands)
func GetDB() *database.DB {
	return db
}

// GetSettings returns the current settings (used by commands)
func GetSettings() *models.Settings {
	return settings
}

// cleanupDatabase closes the database connection
func cleanupDatabase() {
	if db != nil {
		if err := db.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Error closing database: %v\n", err)
		}
	}
}
