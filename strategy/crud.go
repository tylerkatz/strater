package strategy

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylerkatz/strater/config"
)

// Internal API function
func ListStrategiesFromConfig(cfg *config.Config) error {
	fmt.Println("Available Strategies:")
	for _, s := range cfg.Strategies {
		fmt.Printf("- %s", s.Name)
		if s.MonthProfitTargetPct != 0 {
			fmt.Printf(" (Profit Target: %.1f%%)", s.MonthProfitTargetPct*100)
		}
		if s.TradeRiskPct != 0 {
			fmt.Printf(" (Risk: %.1f%%)", s.TradeRiskPct*100)
		}
		fmt.Println()
	}
	return nil
}

// CLI command function
func ListStrategies(cmd *cobra.Command, args []string) error {
	cfgPath, err := config.FindConfigFile()
	if err != nil {
		return err
	}
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
		Name: name,
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
	cfgPath, err := config.FindConfigFile()
	if err != nil {
		return err
	}
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
			// Only update non-zero values
			if profitTarget != 0 {
				cfg.Strategies[i].MonthProfitTargetPct = profitTarget
			}
			if riskPerTrade != 0 {
				cfg.Strategies[i].TradeRiskPct = riskPerTrade
			}
			return nil
		}
	}
	return fmt.Errorf("strategy '%s' not found", name)
}

// AddStrategyCmd handles the CLI command for adding a strategy
func AddStrategyCmd(cmd *cobra.Command, args []string) error {
	cfgPath, err := config.FindConfigFile()
	if err != nil {
		return err
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	if err := AddStrategy(cfg, args[0]); err != nil {
		return err
	}
	return config.SaveConfig(cfg, cfgPath)
}

// RemoveStrategyCmd handles the CLI command for removing a strategy
func RemoveStrategyCmd(cmd *cobra.Command, args []string) error {
	cfgPath, err := config.FindConfigFile()
	if err != nil {
		return err
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	if err := RemoveStrategyFromConfig(cfg, args[0]); err != nil {
		return err
	}
	return config.SaveConfig(cfg, cfgPath)
}

// UpdateStrategyCmd handles the CLI command for updating a strategy
func UpdateStrategyCmd(cmd *cobra.Command, args []string) error {
	cfgPath, err := config.FindConfigFile()
	if err != nil {
		return err
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	profitTarget, _ := cmd.Flags().GetFloat64("profit")
	riskPerTrade, _ := cmd.Flags().GetFloat64("risk")

	if err := UpdateStrategy(cfg, args[0], profitTarget, riskPerTrade); err != nil {
		return err
	}
	return config.SaveConfig(cfg, cfgPath)
}