package strategy

import (
	"fmt"
	"time"

	"github.com/tylerkatz/strater/config"
)

type Transaction struct {
	Date        time.Time
	Type        string // "withdrawal", "reinvestment", "profit"
	Amount      float64
	Balance     float64
	Description string
}

type Account struct {
	Name    string
	Balance float64
	History []Transaction
}

type Plan struct {
	StrategyName   string
	InitialCapital float64
	TradingAccount Account
	SavingsAccount Account
	MonthlyResults []MonthlyResult
}

type MonthlyResult struct {
	Month           int
	StartingBalance float64
	EndingBalance   float64
	ProfitTarget    float64
	RewardPerTrade  float64
	NetWins         int
}

type Analyzer struct {
	config *config.Config
}

func NewAnalyzer(cfg *config.Config) *Analyzer {
	return &Analyzer{config: cfg}
}

func (a *Analyzer) AnalyzeStrategy(strategyName string, months int) (*Plan, error) {
	// Find the requested strategy
	var stratConfig *config.StrategyConfig
	for _, s := range a.config.Strategies {
		if s.Name == strategyName {
			stratConfig = &s
			break
		}
	}
	if stratConfig == nil {
		return nil, fmt.Errorf("strategy not found: %s", strategyName)
	}

	// Initialize the plan with starting capital
	plan := &Plan{
		StrategyName:   strategyName,
		InitialCapital: float64(a.config.Strat.Default.CapitalStart),
		TradingAccount: Account{
			Name:    "Trading",
			Balance: float64(a.config.Strat.Default.CapitalStart),
		},
	}

	// Calculate results for each month
	for month := 1; month <= months; month++ {
		result := a.calculateMonthlyResult(stratConfig, month, plan)
		plan.MonthlyResults = append(plan.MonthlyResults, result)
	}

	return plan, nil
}

func (a *Analyzer) AnalyzeAllStrategies(months int) []*Plan {
	var plans []*Plan
	for _, strat := range a.config.Strategies {
		if plan, err := a.AnalyzeStrategy(strat.Name, months); err == nil {
			plans = append(plans, plan)
		}
	}
	return plans
}

func (a *Analyzer) calculateMonthlyResult(stratConfig *config.StrategyConfig, month int, plan *Plan) MonthlyResult {
	currentBalance := plan.TradingAccount.Balance

	cfg := a.getEffectiveConfig(stratConfig)

	// Calculate monthly profit (reward per trade * number of wins)
	rewardAmount := currentBalance * cfg.TradeRewardPct             // $100 profit per win (1% of 10k)
	monthlyProfit := rewardAmount * float64(cfg.MonthTradesNetWins) // $1000 monthly (10 wins)

	newBalance := currentBalance + monthlyProfit

	plan.TradingAccount.Balance = newBalance

	return MonthlyResult{
		Month:           month,
		StartingBalance: currentBalance,
		EndingBalance:   newBalance,
		ProfitTarget:    monthlyProfit,
		RewardPerTrade:  currentBalance * cfg.TradeRewardPct,
		NetWins:         cfg.MonthTradesNetWins,
	}
}

// getEffectiveConfig returns the merged configuration values, using strategy-specific
// values when present, falling back to defaults when not
func (a *Analyzer) getEffectiveConfig(stratConfig *config.StrategyConfig) struct {
	TradeRewardPct       float64
	TradeRiskPct         float64
	MonthTradesNetWins   int
	MonthProfitTargetPct float64
} {
	// Start with hardcoded defaults
	effective := struct {
		TradeRewardPct       float64
		TradeRiskPct         float64
		MonthTradesNetWins   int
		MonthProfitTargetPct float64
	}{
		TradeRewardPct:       0.01, // 1% default reward
		TradeRiskPct:         0.01, // 1% default risk
		MonthTradesNetWins:   10,   // 10 winning trades default
		MonthProfitTargetPct: 0.10, // 10% monthly profit default
	}

	// Override with config defaults if present
	if a.config.Strat.Default.TradeRewardPct != 0 {
		effective.TradeRewardPct = a.config.Strat.Default.TradeRewardPct
	}
	if a.config.Strat.Default.TradeRiskPct != 0 {
		effective.TradeRiskPct = a.config.Strat.Default.TradeRiskPct
	}
	if a.config.Strat.Default.MonthTradesNetWins != 0 {
		effective.MonthTradesNetWins = a.config.Strat.Default.MonthTradesNetWins
	}
	if a.config.Strat.Default.MonthProfitTargetPct != 0 {
		effective.MonthProfitTargetPct = a.config.Strat.Default.MonthProfitTargetPct
	}

	// Override with strategy-specific values if present
	if stratConfig.TradeRewardPct != 0 {
		effective.TradeRewardPct = stratConfig.TradeRewardPct
	}
	if stratConfig.TradeRiskPct != 0 {
		effective.TradeRiskPct = stratConfig.TradeRiskPct
	}
	if stratConfig.MonthTradesNetWins != 0 {
		effective.MonthTradesNetWins = stratConfig.MonthTradesNetWins
	}
	if stratConfig.MonthProfitTargetPct != 0 {
		effective.MonthProfitTargetPct = stratConfig.MonthProfitTargetPct
	}

	return effective
}

func (a *Analyzer) GetEffectiveConfig(stratConfig *config.StrategyConfig) struct {
	TradeRewardPct       float64
	TradeRiskPct         float64
	MonthTradesNetWins   int
	MonthProfitTargetPct float64
} {
	// Start with hardcoded defaults
	effective := struct {
		TradeRewardPct       float64
		TradeRiskPct         float64
		MonthTradesNetWins   int
		MonthProfitTargetPct float64
	}{
		TradeRewardPct:       0.01, // 1% default reward
		TradeRiskPct:         0.01, // 1% default risk
		MonthTradesNetWins:   10,   // 10 winning trades default
		MonthProfitTargetPct: 0.10, // 10% monthly profit default
	}

	// Override with config defaults if present
	if a.config.Strat.Default.TradeRewardPct != 0 {
		effective.TradeRewardPct = a.config.Strat.Default.TradeRewardPct
	}
	if a.config.Strat.Default.TradeRiskPct != 0 {
		effective.TradeRiskPct = a.config.Strat.Default.TradeRiskPct
	}
	if a.config.Strat.Default.MonthTradesNetWins != 0 {
		effective.MonthTradesNetWins = a.config.Strat.Default.MonthTradesNetWins
	}
	if a.config.Strat.Default.MonthProfitTargetPct != 0 {
		effective.MonthProfitTargetPct = a.config.Strat.Default.MonthProfitTargetPct
	}

	// Override with strategy-specific values if present
	if stratConfig.TradeRewardPct != 0 {
		effective.TradeRewardPct = stratConfig.TradeRewardPct
	}
	if stratConfig.TradeRiskPct != 0 {
		effective.TradeRiskPct = stratConfig.TradeRiskPct
	}
	if stratConfig.MonthTradesNetWins != 0 {
		effective.MonthTradesNetWins = stratConfig.MonthTradesNetWins
	}
	if stratConfig.MonthProfitTargetPct != 0 {
		effective.MonthProfitTargetPct = stratConfig.MonthProfitTargetPct
	}

	return effective
}

// Implementation of calculateMonthlyResult and other helper methods...
