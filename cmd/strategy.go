package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylerkatz/strater/config"
	"github.com/tylerkatz/strater/report"
	"github.com/tylerkatz/strater/strategy"
)

func newStrategyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "strategy",
		Short: "Manage trading strategies",
	}

	analyzeCmd := &cobra.Command{
		Use:   "analyze [strategy-name]",
		Short: "Analyze a trading strategy",
		Args:  cobra.ExactArgs(1),
		RunE:  analyzeStrategy,
	}

	analyzeCmd.Flags().Float64VarP(&initialCapital, "capital", "c", 0, "Initial capital for analysis")
	analyzeCmd.Flags().IntVarP(&months, "months", "m", 12, "Number of months to project")
	analyzeCmd.Flags().StringVarP(&outputFormat, "output", "o", "csv", "Output format (csv, xlsx, or json)")

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List all strategies",
			RunE:  strategy.ListStrategies,
		},
		&cobra.Command{
			Use:   "add [name]",
			Short: "Add a new strategy",
			Args:  cobra.ExactArgs(1),
			RunE:  strategy.AddStrategyCmd,
		},
		&cobra.Command{
			Use:   "remove [name]",
			Short: "Remove a strategy",
			Args:  cobra.ExactArgs(1),
			RunE:  strategy.RemoveStrategyCmd,
		},
		&cobra.Command{
			Use:   "update [name]",
			Short: "Update a strategy",
			Args:  cobra.ExactArgs(1),
			RunE:  strategy.UpdateStrategyCmd,
		},
		analyzeCmd,
	)

	return cmd
}

func analyzeStrategy(cmd *cobra.Command, args []string) error {
	cfgPath := configPath
	if cfgPath == "" {
		cfgPath = config.FindConfigFile()
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	analyzer := strategy.NewAnalyzer(cfg)
	plan, err := analyzer.AnalyzeStrategy(args[0], months)
	if err != nil {
		return fmt.Errorf("error analyzing strategy %s: %v", args[0], err)
	}

	outputPath := fmt.Sprintf("%s_analysis.%s", args[0], outputFormat)
	if err := report.Generate([]*strategy.Plan{plan}, outputFormat, outputPath); err != nil {
		return fmt.Errorf("error generating report: %v", err)
	}

	fmt.Printf("Analysis complete: %s\n", outputPath)
	return nil
}
