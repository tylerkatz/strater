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
	RiskPerTrade    float64
	Withdrawal      float64
	Transactions    []Transaction
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

	plan := &Plan{
		StrategyName:   strategyName,
		InitialCapital: float64(a.config.Strat.Default.CapitalStart),
		TradingAccount: Account{Name: "Trading", Balance: float64(a.config.Strat.Default.CapitalStart)},
		SavingsAccount: Account{Name: "Savings", Balance: 0},
	}

	// Calculate monthly results
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
	profitTarget := currentBalance * stratConfig.MonthlyProfitTarget
	riskPerTrade := currentBalance * stratConfig.RiskPerTrade

	return MonthlyResult{
		Month:           month,
		StartingBalance: currentBalance,
		EndingBalance:   currentBalance + profitTarget,
		ProfitTarget:    profitTarget,
		RiskPerTrade:    riskPerTrade,
	}
}

// Implementation of calculateMonthlyResult and other helper methods...
