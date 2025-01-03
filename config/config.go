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
		"config.json",                       // Current directory
		"$HOME/.config/strater/config.json", // User's config directory
		"/etc/strater/config.json",          // System-wide config
	}

	// Can be overridden by STRATER_CONFIG env var
	ConfigEnvVar = "STRATER_CONFIG"
)

type Config struct {
	Strat struct {
		Default struct {
			CapitalStart      int     `json:"capital_start"`
			TradeRisk         float64 `json:"trade_risk"`
			MonthProfitTarget float64 `json:"month_profit_target"`
			MonthCount        int     `json:"month_count"`
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
	Name                string  `json:"name"`
	MonthlyProfitTarget float64 `json:"monthly_profit_target"`
	RiskPerTrade        float64 `json:"risk_per_trade"`
}

// FindConfigFile looks for config file in standard locations
func FindConfigFile() string {
	// First check environment variable
	if envPath := os.Getenv(ConfigEnvVar); envPath != "" {
		return envPath
	}

	// Then check standard locations
	for _, path := range DefaultConfigPaths {
		// Expand environment variables like $HOME
		expandedPath := os.ExpandEnv(path)
		if _, err := os.Stat(expandedPath); err == nil {
			return expandedPath
		}
	}

	// Default to local config if nothing found
	return DefaultConfigPaths[0]
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return &Config{
				Strat: struct {
					Default struct {
						CapitalStart      int     `json:"capital_start"`
						TradeRisk         float64 `json:"trade_risk"`
						MonthProfitTarget float64 `json:"month_profit_target"`
						MonthCount        int     `json:"month_count"`
					} `json:"default"`
				}{
					Default: struct {
						CapitalStart      int     `json:"capital_start"`
						TradeRisk         float64 `json:"trade_risk"`
						MonthProfitTarget float64 `json:"month_profit_target"`
						MonthCount        int     `json:"month_count"`
					}{
						CapitalStart:      10000,
						TradeRisk:         0.02,
						MonthProfitTarget: 0.20,
						MonthCount:        12,
					},
				},
				Strategies: []StrategyConfig{},
			}, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
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
	case KeyCapitalStart:
		fmt.Printf("%d\n", cfg.Strat.Default.CapitalStart)
	case KeyTradeRisk:
		fmt.Printf("%.2f\n", cfg.Strat.Default.TradeRisk)
	case KeyMonthProfitTarget:
		fmt.Printf("%.2f\n", cfg.Strat.Default.MonthProfitTarget)
	case KeyMonthCount:
		fmt.Printf("%d\n", cfg.Strat.Default.MonthCount)
	default:
		return fmt.Errorf("unknown config key: %s\nAvailable keys: %v", path, GetAvailableKeys())
	}
	return nil
}

func SetConfigValue(cfg *Config, path string, value string) error {
	switch path {
	case KeyCapitalStart:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s - must be a number", path, value)
		}
		cfg.Strat.Default.CapitalStart = int(v)
	case KeyTradeRisk:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s", path, value)
		}
		if v > 1.0 {
			return fmt.Errorf("trade risk cannot be more than 1 (100%%)")
		}
		cfg.Strat.Default.TradeRisk = v
	case KeyMonthProfitTarget:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s", path, value)
		}
		if v <= 0 {
			return fmt.Errorf("month profit target must be positive")
		}
		cfg.Strat.Default.MonthProfitTarget = v
	case KeyMonthCount:
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %s - must be an integer", path, value)
		}
		if v <= 0 {
			return fmt.Errorf("month count must be positive")
		}
		cfg.Strat.Default.MonthCount = v
	default:
		return fmt.Errorf("unknown config key: %s\nAvailable keys: %v", path, GetAvailableKeys())
	}
	return nil
}
