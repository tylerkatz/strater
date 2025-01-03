package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tylerkatz/strater/config"
	"github.com/tylerkatz/strater/report"
	"github.com/tylerkatz/strater/strategy"
)

func newStrategyCmd() *cobra.Command {
	var outputPath string
	var listKeys bool

	cmd := &cobra.Command{
		Use:     "strat",
		Aliases: []string{"strategy"},
		Short:   "Manage trading strategies",
	}

	configCmd := &cobra.Command{
		Use:   "config [strategy-name] [key] [value]",
		Short: "Get or set strategy configuration values",
		Long: `Get or set strategy-specific configuration values.
Examples: 
  strater strat config ifunds trade_reward_pct 0.02    # Set a value
  strater strat config ifunds trade_reward_pct         # Get a value
  strater strat config ifunds -l                       # List available keys`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if listKeys {
				if len(args) < 1 {
					return fmt.Errorf("strategy name is required")
				}
				strategyName := args[0]
				cfgPath := configPath
				if cfgPath == "" {
					var err error
					cfgPath, err = config.FindConfigFile()
					if err != nil {
						return err
					}
				}

				cfg, err := config.LoadConfig(cfgPath)
				if err != nil {
					return fmt.Errorf("error loading config: %v", err)
				}

				// Find the strategy
				var strat config.StrategyConfig
				found := false
				for _, s := range cfg.Strategies {
					if s.Name == strategyName {
						strat = s
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("strategy '%s' not found", strategyName)
				}

				// Get effective config
				analyzer := strategy.NewAnalyzer(cfg)
				effective := analyzer.GetEffectiveConfig(&strat)

				fmt.Printf("%s: %s\n", config.KeyStratName, strat.Name)
				fmt.Printf("%s: %s\n", config.KeyStratDescription, strat.Description)
				fmt.Printf("%s: %.2f\n", config.KeyStratTradeRiskPct, effective.TradeRiskPct)
				fmt.Printf("%s: %.2f\n", config.KeyStratTradeRewardPct, effective.TradeRewardPct)
				fmt.Printf("%s: %d\n", config.KeyStratMonthTradesNetWins, effective.MonthTradesNetWins)
				fmt.Printf("%s: %.2f\n", config.KeyStratMonthProfitTargetPct, effective.MonthProfitTargetPct)
				return nil
			}
			if len(args) < 2 {
				return fmt.Errorf("requires at least strategy name and key: strater strat config <strategy> <key> [value]")
			}
			return handleStrategyConfig(cmd, args)
		},
		Args: cobra.ArbitraryArgs,
	}

	configCmd.Flags().BoolVarP(&listKeys, "list", "l", false, "List available configuration keys")

	analyzeCmd := &cobra.Command{
		Use:   "analyze [strategy-name]",
		Short: "Analyze a trading strategy",
		Long: `Analyze a trading strategy using default values from config or override with flags.
Example: strater strategy analyze conservative

All flags are optional and will use config defaults if not specified.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return analyzeStrategy(cmd, args, outputPath)
		},
	}

	analyzeCmd.Flags().Float64VarP(&initialCapital, "capital", "c", 0, "Initial capital (default from config)")
	analyzeCmd.Flags().IntVarP(&months, "months", "m", 0, "Number of months to project (default from config)")
	analyzeCmd.Flags().StringVarP(&outputFormat, "output", "o", "csv", "Output format (csv, xlsx, or json)")
	analyzeCmd.Flags().StringVarP(&outputPath, "path", "p", "", "Output file path (default from config)")

	// Add existing commands
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
		configCmd,
		analyzeCmd,
	)

	return cmd
}

func handleStrategyConfig(_ *cobra.Command, args []string) error {
	cfgPath, err := config.FindConfigFile()
	if err != nil {
		return err
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	strategyName := args[0]
	key := args[1]

	// Find the strategy
	var stratIndex = -1
	for i, s := range cfg.Strategies {
		if s.Name == strategyName {
			stratIndex = i
			break
		}
	}
	if stratIndex == -1 {
		return fmt.Errorf("strategy '%s' not found", strategyName)
	}

	// If no value provided, get the current value
	if len(args) == 2 {
		switch key {
		case config.KeyStratName:
			fmt.Printf("%s\n", cfg.Strategies[stratIndex].Name)
		case config.KeyStratDescription:
			fmt.Printf("%s\n", cfg.Strategies[stratIndex].Description)
		case config.KeyStratTradeRiskPct:
			fmt.Printf("%.2f\n", cfg.Strategies[stratIndex].TradeRiskPct)
		case config.KeyStratTradeRewardPct:
			fmt.Printf("%.2f\n", cfg.Strategies[stratIndex].TradeRewardPct)
		case config.KeyStratMonthTradesNetWins:
			fmt.Printf("%d\n", cfg.Strategies[stratIndex].MonthTradesNetWins)
		case config.KeyStratMonthProfitTargetPct:
			fmt.Printf("%.2f\n", cfg.Strategies[stratIndex].MonthProfitTargetPct)
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}
		return nil
	}

	// Set the new value
	value := args[2]
	switch key {
	case config.KeyStratName:
		cfg.Strategies[stratIndex].Name = value
	case config.KeyStratDescription:
		cfg.Strategies[stratIndex].Description = value
	case config.KeyStratTradeRiskPct, config.KeyStratTradeRewardPct, config.KeyStratMonthProfitTargetPct:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s - must be a number", key, value)
		}
		if v <= 0 {
			return fmt.Errorf("%s must be positive", key)
		}
		switch key {
		case config.KeyStratTradeRiskPct:
			cfg.Strategies[stratIndex].TradeRiskPct = v
		case config.KeyStratTradeRewardPct:
			cfg.Strategies[stratIndex].TradeRewardPct = v
		case config.KeyStratMonthProfitTargetPct:
			cfg.Strategies[stratIndex].MonthProfitTargetPct = v
		}
	case config.KeyStratMonthTradesNetWins:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s - must be an integer", key, value)
		}
		cfg.Strategies[stratIndex].MonthTradesNetWins = v
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	return config.SaveConfig(cfg, cfgPath)
}

func analyzeStrategy(_ *cobra.Command, args []string, outputPath string) error {
	cfgPath := configPath
	if cfgPath == "" {
		var err error
		cfgPath, err = config.FindConfigFile()
		if err != nil {
			return err
		}
	}

	// Debug config loading
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	// Use config defaults if flags aren't set
	if months == 0 {
		months = cfg.Strat.Default.MonthCount
	}
	if initialCapital == 0 {
		initialCapital = float64(cfg.Strat.Default.CapitalStart)
	}

	analyzer := strategy.NewAnalyzer(cfg)
	plan, err := analyzer.AnalyzeStrategy(args[0], months)
	if err != nil {
		return fmt.Errorf("error analyzing strategy %s: %v", args[0], err)
	}

	// Use config output path if not overridden
	if outputPath == "" && cfg.Settings.OutputPath != "" {
		outputPath = fmt.Sprintf("%s/%s_analysis.%s",
			cfg.Settings.OutputPath,
			args[0],
			outputFormat)
	} else {
		outputPath = fmt.Sprintf("%s_analysis.%s", args[0], outputFormat)
	}

	if err := report.Generate([]*strategy.Plan{plan}, outputFormat, outputPath); err != nil {
		return fmt.Errorf("error generating report: %v", err)
	}

	return nil
}
