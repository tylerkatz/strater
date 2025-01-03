package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configPath     string
	outputFormat   string
	months         int
	initialCapital float64
)

var rootCmd = &cobra.Command{
	Use:   "strater",
	Short: "A trading strategy scaling calculator",
	Long:  `Strater helps you plan and analyze account and risk management strategies.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add all subcommands
	rootCmd.AddCommand(newStrategyCmd())
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newInitCmd())
}
