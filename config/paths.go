package config

// Config key paths
const (
	KeyCapitalStart      = "strat.default.capital_start"
	KeyTradeRisk         = "strat.default.trade_risk"
	KeyMonthProfitTarget = "strat.default.month_profit_target"
	KeyMonthCount        = "strat.default.month_count"
)

// GetAvailableKeys returns a list of all available config keys
func GetAvailableKeys() []string {
	return []string{
		KeyCapitalStart,
		KeyTradeRisk,
		KeyMonthProfitTarget,
		KeyMonthCount,
	}
}
