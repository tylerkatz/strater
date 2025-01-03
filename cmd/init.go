package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tylerkatz/strater/config"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new configuration file",
		Long:  `Creates a new .strater configuration file in the current directory with default values.`,
		RunE:  handleInit,
	}
	return cmd
}

func handleInit(_ *cobra.Command, _ []string) error {
	// Check if config file already exists in current directory
	if _, err := os.Stat(".strater.json"); err == nil {
		return fmt.Errorf("a .strater.json configuration file already exists in the current directory")
	}

	// Check if config exists in other locations
	if path, err := config.FindConfigFile(); err == nil {
		return fmt.Errorf("configuration file already exists at %s", path)
	}

	// Create default config
	cfg := &config.Config{
		Settings: struct {
			OutputPath string `json:"output_path"`
		}{
			OutputPath: "strater_output",
		},
		Strat: struct {
			Default struct {
				CapitalStart         int     `json:"capital_start"`
				TradeRiskPct         float64 `json:"trade_risk_pct"`
				TradeRewardPct       float64 `json:"trade_reward_pct"`
				MonthTradesNetWins   int     `json:"month_trades_net_wins"`
				MonthProfitTargetPct float64 `json:"month_profit_target_pct"`
				MonthCount           int     `json:"month_count"`
			} `json:"default"`
		}{
			Default: struct {
				CapitalStart         int     `json:"capital_start"`
				TradeRiskPct         float64 `json:"trade_risk_pct"`
				TradeRewardPct       float64 `json:"trade_reward_pct"`
				MonthTradesNetWins   int     `json:"month_trades_net_wins"`
				MonthProfitTargetPct float64 `json:"month_profit_target_pct"`
				MonthCount           int     `json:"month_count"`
			}{
				CapitalStart:         10000,
				TradeRiskPct:         0.01,
				TradeRewardPct:       0.01,
				MonthTradesNetWins:   10,
				MonthProfitTargetPct: 0.10,
				MonthCount:           12,
			},
		},
		Strategies: []config.StrategyConfig{},
	}

	// Save the config file
	if err := config.SaveConfig(cfg, ".strater.json"); err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}

	fmt.Println("Created new configuration file: .strater.json")
	return nil
}
