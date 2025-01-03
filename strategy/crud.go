package strategy

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/projectname/config"
)

// Internal API function
func ListStrategiesFromConfig(cfg *config.Config) error {
	fmt.Println("Available Strategies:")
	for _, s := range cfg.Strategies {
		fmt.Printf("- %s (Profit Target: %.1f%%, Risk: %.1f%%)\n",
			s.Name,
			s.MonthlyProfitTarget*100,
			s.RiskPerTrade*100)
	}
	return nil
}

// CLI command function
func ListStrategies(cmd *cobra.Command, args []string) error {
	cfgPath := config.FindConfigFile()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}
	return ListStrategiesFromConfig(cfg)
}

// AddStrategy adds a new strategy to the config
func AddStrategy(cfg *config.Config, name string) error {
	for _, s := range cfg.Strategies {
		if s.Name == name {
			return fmt.Errorf("strategy '%s' already exists", name)
		}
	}

	newStrategy := config.StrategyConfig{
		Name:                name,
		MonthlyProfitTarget: 0.20, // Default 20%
		RiskPerTrade:        0.02, // Default 2%
	}

	cfg.Strategies = append(cfg.Strategies, newStrategy)
	return nil
}

// Internal API function
func RemoveStrategyFromConfig(cfg *config.Config, name string) error {
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
	return nil
}

// CLI command function
func RemoveStrategy(cmd *cobra.Command, args []string) error {
	cfgPath := config.FindConfigFile()
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	if err := RemoveStrategyFromConfig(cfg, args[0]); err != nil {
		return err
	}
	return config.SaveConfig(cfg, cfgPath)
}

// UpdateStrategy updates an existing strategy
func UpdateStrategy(cfg *config.Config, name string, profitTarget, riskPerTrade float64) error {
	for i, s := range cfg.Strategies {
		if s.Name == name {
			cfg.Strategies[i].MonthlyProfitTarget = profitTarget
			cfg.Strategies[i].RiskPerTrade = riskPerTrade
			return nil
		}
	}
	return fmt.Errorf("strategy '%s' not found", name)
}
