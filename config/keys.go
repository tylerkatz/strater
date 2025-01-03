package config

// Strategy configuration keys - shared between defaults and strategy-specific settings
const (
	// Strategy metadata
	KeyStratName        = "name"
	KeyStratDescription = "description"

	// Core trading parameters
	KeyStratTradeRiskPct         = "trade_risk_pct"
	KeyStratTradeRewardPct       = "trade_reward_pct"
	KeyStratMonthTradesNetWins   = "month_trades_net_wins"
	KeyStratMonthProfitTargetPct = "month_profit_target_pct"

	// Default-only parameters
	KeyStratDefaultCapitalStart = "capital_start"
	KeyStratDefaultMonthCount   = "month_count"
	KeySettingsOutputPath       = "output_path"
)

// GetStrategyKeys returns keys that can be used in strategy-specific configs
func GetStrategyKeys() []string {
	return []string{
		KeyStratName,
		KeyStratDescription,
		KeyStratTradeRiskPct,
		KeyStratTradeRewardPct,
		KeyStratMonthTradesNetWins,
		KeyStratMonthProfitTargetPct,
	}
}

// GetDefaultKeys returns all available config keys for defaults
func GetDefaultKeys() []string {
	return []string{
		KeyStratDefaultCapitalStart,
		KeyStratTradeRiskPct,
		KeyStratTradeRewardPct,
		KeyStratMonthTradesNetWins,
		KeyStratMonthProfitTargetPct,
		KeyStratDefaultMonthCount,
		KeySettingsOutputPath,
	}
}
