package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tylerkatz/strater/config"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config [key] [value]",
		Short: "Get or set configuration values",
		Long:  `Get or set configuration values using dot notation (e.g., strat.default.trade_risk)`,
		Args:  cobra.RangeArgs(1, 2),
		RunE:  handleConfig,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all available configuration keys",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			for _, key := range config.GetAvailableKeys() {
				fmt.Printf("%s\n", key)
			}
			return nil
		},
	})

	cmd.PersistentFlags().StringVarP(&configPath, "file", "f", "", "Override default config file location")
	return cmd
}

func handleConfig(cmd *cobra.Command, args []string) error {
	cfgPath := configPath
	if cfgPath == "" {
		cfgPath = config.FindConfigFile()
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	key := args[0]
	if len(args) == 2 {
		if err := config.SetConfigValue(cfg, key, args[1]); err != nil {
			return err
		}
		return config.SaveConfig(cfg, cfgPath)
	}

	return config.GetConfigValue(cfg, key)
}
