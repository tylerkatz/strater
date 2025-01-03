package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/projectname/config"
)

var (
	profitTarget float64
	riskPerTrade float64
)

func init() {
	updateCmd := &cobra.Command{
		Use:   "update [name]",
		Short: "Update a strategy",
		Args:  cobra.ExactArgs(1),
		RunE:  updateStrategy,
	}

	updateCmd.Flags().Float64VarP(&profitTarget, "profit", "p", 0.20, "Monthly profit target (as decimal)")
	updateCmd.Flags().Float64VarP(&riskPerTrade, "risk", "r", 0.02, "Risk per trade (as decimal)")
}

func addStrategy(_ *cobra.Command, args []string) error {
	cfgPath := config.FindConfigFile()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	name := args[0]
	// Check if strategy already exists
	for _, s := range cfg.Strategies {
		if s.Name == name {
			return fmt.Errorf("strategy '%s' already exists", name)
		}
	}

	// Add new strategy
	newStrategy := config.StrategyConfig{
		Name:                name,
		MonthlyProfitTarget: 0.20, // Default 20%
		RiskPerTrade:        0.02, // Default 2%
	}

	cfg.Strategies = append(cfg.Strategies, newStrategy)
	return config.SaveConfig(cfg, cfgPath)
}

func removeStrategy(_ *cobra.Command, args []string) error {
	cfgPath := config.FindConfigFile()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	name := args[0]
	found := false
	for i, s := range cfg.Strategies {
		if s.Name == name {
			cfg.Strategies = append(cfg.Strategies[:i], cfg.Strategies[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("strategy '%s' not found", name)
	}

	return config.SaveConfig(cfg, cfgPath)
}

func updateStrategy(cmd *cobra.Command, args []string) error {
	cfgPath := config.FindConfigFile()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	name := args[0]
	for i, s := range cfg.Strategies {
		if s.Name == name {
			// Update strategy values
			cfg.Strategies[i].MonthlyProfitTarget = profitTarget
			cfg.Strategies[i].RiskPerTrade = riskPerTrade
			return config.SaveConfig(cfg, cfgPath)
		}
	}

	return fmt.Errorf("strategy '%s' not found", name)
}
