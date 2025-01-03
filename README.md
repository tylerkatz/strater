# Strater

Strater is a command-line tool that helps traders plan their account progression and capital growth through multiple stages. It allows you to model different risk management and capital allocation strategies across accounts, helping you chart a clear path from your starting capital to your desired portfolio.

Whether you're planning to scale from a small account to larger ones, managing multiple trading accounts with different risk parameters, or developing a progression strategy for fund management, Strater helps you:

- Model different account growth scenarios
- Plan risk-appropriate position sizing at each stage
- Track progress against your scaling targets
- Analyze potential outcomes across multiple accounts
- Export detailed progression plans and analysis

## Features

### Core Features (✅ Implemented)
- Strategy Management
  - Create and manage multiple trading strategies
  - Configure risk and reward parameters
  - Set monthly targets and goals
- Analysis Output
  - CSV export
  - XLSX (Excel) export
  - JSON export
- Configuration System
  - Default settings
  - Per-strategy customization
  - Multiple config file locations

### Under Development (🚧 In Progress)
1. **Portfolio Account Management**
   - Multi-account strategy modeling
   - Separate investment vs equity tracking
   - Account types:
     - Managed accounts (low investment, high equity)
     - Owned accounts (high investment, high equity)
   - Account lifecycle events:
     - Account opening/closing
     - Investment deposits/withdrawals
     - Equity changes
     - Balance transfers between accounts

2. **Transaction Ledger**
   - Comprehensive event tracking
   - Management transaction history
   - Account lifecycle events:
     - Account creation/closure records
     - Investment vs equity tracking
     - Balance transfers between accounts
   - Investment flow tracking
   - Balance and equity reconciliation

3. **Account Progression Modeling**
   - Path from managed to owned accounts
   - Capital allocation strategies
   - Investment requirement planning
   - Equity scaling projections
   - Portfolio transition timing

### Planned Features (📋 Upcoming)
1. **Position Sizing Calculator**
   - Calculate optimal position sizes based on risk parameters
   - Support for different position sizing methods
   - Risk-based lot size calculations

2. **Risk Metrics**
   - Maximum drawdown calculations
   - Sharpe ratio analysis
   - Risk-adjusted return metrics
   - Win/loss ratio tracking

3. **Portfolio Analysis**
   - Multi-account performance metrics
   - Investment vs equity ratios
   - Capital efficiency calculations
   - Portfolio composition tracking

4. **Report Enhancements**
   - Custom report templates
   - PDF generation
   - Interactive visualizations
   - Transaction audit trails

### Future Enhancements (🔮 Planned)
1. **Monte Carlo Simulation**
   - Randomized outcome modeling
   - Probability distribution analysis
   - Scenario visualization

2. **Web Interface**
   - Browser-based UI
   - Interactive strategy builder
   - Real-time monitoring

3. **API Support**
   - REST API endpoints
   - Broker API integration
   - Strategy webhooks

## Installation

### Option 1: Install from Source (Recommended)
```bash
go install github.com/tylerkatz/strater@latest
```

This will install the `strater` binary to your `$GOPATH/bin` directory. Make sure your `$GOPATH/bin` is in your system's PATH to run `strater` commands from anywhere.

### Option 2: Build from Source
1. Clone the repository:
```bash
git clone https://github.com/tylerkatz/strater.git
cd strater
```

2. Install dependencies and build:
```bash
go mod download
go build
```

3. Run the binary:
```bash
./strater
```

## Quick Start

1. Initialize a new configuration:
```bash
strater init
```

2. Add a trading strategy:
```bash
strater strat add conservative
```

3. Configure your strategy:
```bash
strater strat config conservative trade_risk_pct 0.01
strater strat config conservative trade_reward_pct 0.02
```

4. Analyze the strategy:
```bash
strater strat analyze conservative --months 12 --output xlsx
```

## Configuration

Strater looks for configuration files in the following locations:
- `.strater.json` in the current directory
- `$HOME/.config/strater/.strater.json`
- `/etc/strater/.strater.json`

You can also specify a custom config location using the `STRATER_CONFIG` environment variable.

### Available Configuration Keys

- `capital_start`: Initial trading capital
- `trade_risk_pct`: Risk percentage per trade
- `trade_reward_pct`: Reward percentage per trade
- `month_trades_net_wins`: Target number of winning trades per month
- `month_profit_target_pct`: Monthly profit target percentage
- `month_count`: Number of months to analyze
- `output_path`: Default path for analysis output files

## Usage Examples

List all strategies:
```bash
strater strat list
```

View strategy configuration:
```bash
strater strat config mystrategy -l
```

Analyze with custom parameters:
```bash
strater strat analyze mystrategy --capital 25000 --months 24 --output json
```

## Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please make sure to update tests as appropriate and adhere to the existing coding style.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to all contributors who have helped shape Strater
- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Uses [excelize](https://github.com/xuri/excelize) for Excel file generation

## Support

If you encounter any problems or have suggestions, please [open an issue](https://github.com/tylerkatz/strater/issues).

## Security

For security issues, please open an issue but do not disclose any specific details.  You will be provided with a link where you can disclose the issue.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/tylerkatz/strater/tags).

## Authors

* **Tyler Katz** - *Initial work*

See also the list of [contributors](https://github.com/tylerkatz/strater/contributors) who participated in this project.
