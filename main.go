package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yourusername/projectname/config"
	"github.com/yourusername/projectname/report"
	"github.com/yourusername/projectname/strategy"
)

var (
	// CLI flags
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

func init() {
	// Strategy management commands
	strategyCmd := &cobra.Command{
		Use:   "strategy",
		Short: "Manage trading strategies",
	}

	// Analysis command
	analyzeCmd := &cobra.Command{
		Use:   "analyze [strategy-name]",
		Short: "Analyze a trading strategy",
		Args:  cobra.ExactArgs(1),
		RunE:  analyzeStrategy,
	}

	// Add flags to analyze command
	analyzeCmd.Flags().Float64VarP(&initialCapital, "capital", "c", 0, "Initial capital for analysis")
	analyzeCmd.Flags().IntVarP(&months, "months", "m", 12, "Number of months to project")
	analyzeCmd.Flags().StringVarP(&outputFormat, "output", "o", "csv", "Output format (csv, xlsx, or json)")

	strategyCmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List all strategies",
			RunE:  listStrategies,
		},
		&cobra.Command{
			Use:   "add [name]",
			Short: "Add a new strategy",
			Args:  cobra.ExactArgs(1),
			RunE:  addStrategy,
		},
		&cobra.Command{
			Use:   "remove [name]",
			Short: "Remove a strategy",
			Args:  cobra.ExactArgs(1),
			RunE:  removeStrategy,
		},
		&cobra.Command{
			Use:   "update [name]",
			Short: "Update a strategy",
			Args:  cobra.ExactArgs(1),
			RunE:  updateStrategy,
		},
		analyzeCmd,
	)

	// Config management commands
	configCmd := &cobra.Command{
		Use:   "config [key] [value]",
		Short: "Get or set configuration values",
		Long:  `Get or set configuration values using dot notation (e.g., strat.default.trade_risk)`,
		Args:  cobra.RangeArgs(1, 2),
		RunE:  handleConfig,
	}

	configCmd.PersistentFlags().StringVarP(&configPath, "file", "f", "", "Override default config file location")
	rootCmd.AddCommand(configCmd)
}

func listStrategies(_ *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}
	return strategy.ListStrategiesFromConfig(cfg)
}

func addStrategy(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	if err := strategy.AddStrategy(cfg, args[0]); err != nil {
		return err
	}
	return config.SaveConfig(cfg, configPath)
}

func removeStrategy(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	if err := strategy.RemoveStrategyFromConfig(cfg, args[0]); err != nil {
		return err
	}
	return config.SaveConfig(cfg, configPath)
}

func updateStrategy(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	if err := strategy.UpdateStrategy(cfg, args[0], 0.20, 0.02); err != nil {
		return err
	}
	return config.SaveConfig(cfg, configPath)
}

func analyzeStrategy(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	strategyName := args[0]
	analyzer := strategy.NewAnalyzer(cfg)

	plan, err := analyzer.AnalyzeStrategy(strategyName, months)
	if err != nil {
		return fmt.Errorf("error analyzing strategy %s: %v", strategyName, err)
	}

	outputPath := fmt.Sprintf("%s_analysis.%s", strategyName, outputFormat)
	if err := report.Generate([]*strategy.Plan{plan}, outputFormat, outputPath); err != nil {
		return fmt.Errorf("error generating report: %v", err)
	}

	fmt.Printf("Analysis complete: %s\n", outputPath)
	return nil
}

func handleConfig(cmd *cobra.Command, args []string) error {
	// Use provided path or find default
	cfgPath := configPath
	if cfgPath == "" {
		cfgPath = config.FindConfigFile()
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	key := args[0]
	// If a value is provided, set the config
	if len(args) == 2 {
		value := args[1]
		if err := setConfigValue(cfg, key, value); err != nil {
			return err
		}
		return config.SaveConfig(cfg, cfgPath)
	}

	// Otherwise, get the config value
	return getConfigValue(cfg, key)
}

func getConfigValue(cfg *config.Config, key string) error {
	switch key {
	case "strat.default.capital_start":
		fmt.Printf("%d\n", cfg.Strat.Default.CapitalStart)
	case "strat.default.trade_risk":
		fmt.Printf("%.2f\n", cfg.Strat.Default.TradeRisk)
	case "strat.default.month_profit_target":
		fmt.Printf("%.2f\n", cfg.Strat.Default.MonthProfitTarget)
	case "strat.default.month_count":
		fmt.Printf("%d\n", cfg.Strat.Default.MonthCount)
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}

func setConfigValue(cfg *config.Config, key string, value string) error {
	switch key {
	case "strat.default.capital_start":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		cfg.Strat.Default.CapitalStart = int(v)
	case "strat.default.trade_risk":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		cfg.Strat.Default.TradeRisk = v
	case "strat.default.month_profit_target":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		cfg.Strat.Default.MonthProfitTarget = v
	case "strat.default.month_count":
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s", key, value)
		}
		cfg.Strat.Default.MonthCount = v
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
