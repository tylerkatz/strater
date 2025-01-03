package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var (
	// DefaultConfigPaths lists paths to check in order
	DefaultConfigPaths = []string{
		".strater.json",                       // Current directory
		"$HOME/.config/strater/.strater.json", // User's config directory
		"/etc/strater/.strater.json",          // System-wide config
	}

	// Can be overridden by STRATER_CONFIG env var
	ConfigEnvVar = "STRATER_CONFIG"
)

type Config struct {
	Settings struct {
		OutputPath string `json:"output_path"`
	} `json:"settings"`
	Strat struct {
		Default struct {
			CapitalStart         int     `json:"capital_start"`
			TradeRiskPct         float64 `json:"trade_risk_pct"`
			TradeRewardPct       float64 `json:"trade_reward_pct"`
			MonthTradesNetWins   int     `json:"month_trades_net_wins"`
			MonthProfitTargetPct float64 `json:"month_profit_target_pct"`
			MonthCount           int     `json:"month_count"`
		} `json:"default"`
	} `json:"strat"`
	Strategies []StrategyConfig `json:"strategies"`
}

type DefaultConfig struct {
	InitialCapital      int     `json:"initial_capital"` // e.g., 10000
	DefaultRiskPerTrade float64 `json:"default_risk_per_trade"`
	DefaultProfitTarget float64 `json:"default_profit_target"`
	DefaultMonths       int     `json:"default_months"`
}

type StrategyConfig struct {
	Name                 string  `json:"name"`
	Description          string  `json:"description,omitempty"`
	TradeRiskPct         float64 `json:"trade_risk_pct,omitempty"`
	TradeRewardPct       float64 `json:"trade_reward_pct,omitempty"`
	MonthTradesNetWins   int     `json:"month_trades_net_wins,omitempty"`
	MonthProfitTargetPct float64 `json:"month_profit_target_pct,omitempty"`
}

// FindConfigFile looks for config file in standard locations
func FindConfigFile() (string, error) {
	// First check environment variable
	if envPath := os.Getenv(ConfigEnvVar); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
	}

	// Then check standard locations
	for _, path := range DefaultConfigPaths {
		// Expand environment variables like $HOME
		expandedPath := os.ExpandEnv(path)
		if _, err := os.Stat(expandedPath); err == nil {
			return expandedPath, nil
		}
	}

	return "", fmt.Errorf("no configuration file found. Run 'strater init' to create one")
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return getDefaultConfig(), nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Helper function to return default config
func getDefaultConfig() *Config {
	return &Config{
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
		Strategies: []StrategyConfig{},
	}
}

func SaveConfig(cfg *Config, path string) error {
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func GetConfigValue(cfg *Config, path string) error {
	switch path {
	case KeyStratDefaultCapitalStart:
		fmt.Printf("%d\n", cfg.Strat.Default.CapitalStart)
	case KeyStratTradeRiskPct:
		fmt.Printf("%.2f\n", cfg.Strat.Default.TradeRiskPct)
	case KeyStratMonthProfitTargetPct:
		fmt.Printf("%.2f\n", cfg.Strat.Default.MonthProfitTargetPct)
	case KeyStratDefaultMonthCount:
		fmt.Printf("%d\n", cfg.Strat.Default.MonthCount)
	case KeyStratTradeRewardPct:
		fmt.Printf("%.2f\n", cfg.Strat.Default.TradeRewardPct)
	case KeyStratMonthTradesNetWins:
		fmt.Printf("%d\n", cfg.Strat.Default.MonthTradesNetWins)
	case KeySettingsOutputPath:
		fmt.Printf("%s\n", cfg.Settings.OutputPath)
	default:
		return fmt.Errorf("unknown config key: %s\nAvailable keys: %v", path, GetDefaultKeys())
	}
	return nil
}

func SetConfigValue(cfg *Config, path string, value string) error {
	switch path {
	case KeyStratDefaultCapitalStart:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s - must be a number", path, value)
		}
		cfg.Strat.Default.CapitalStart = int(v)
	case KeyStratTradeRiskPct:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s", path, value)
		}
		if v > 1.0 {
			return fmt.Errorf("trade risk cannot be more than 1 (100%%)")
		}
		cfg.Strat.Default.TradeRiskPct = v
	case KeyStratMonthProfitTargetPct:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s", path, value)
		}
		if v <= 0 {
			return fmt.Errorf("month profit target must be positive")
		}
		cfg.Strat.Default.MonthProfitTargetPct = v
	case KeyStratDefaultMonthCount:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s - must be an integer", path, value)
		}
		if v <= 0 {
			return fmt.Errorf("month count must be positive")
		}
		cfg.Strat.Default.MonthCount = v
	case KeyStratTradeRewardPct:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s - must be a number", path, value)
		}
		if v <= 0 {
			return fmt.Errorf("trade reward must be positive")
		}
		cfg.Strat.Default.TradeRewardPct = v
	case KeyStratMonthTradesNetWins:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s - must be an integer", path, value)
		}
		cfg.Strat.Default.MonthTradesNetWins = v
	case KeySettingsOutputPath:
		cfg.Settings.OutputPath = value
	default:
		return fmt.Errorf("unknown config key: %s\nAvailable keys: %v", path, GetDefaultKeys())
	}
	return nil
}
