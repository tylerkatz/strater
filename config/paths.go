package config

// Config path prefixes
const (
	SettingsPrefix = "settings."
	StratPrefix    = "strat."
	DefaultPrefix  = StratPrefix + "default."
)

// GetDefaultPath returns the full config path for a default setting
func GetDefaultPath(key string) string {
	if key == KeySettingsOutputPath {
		return SettingsPrefix + key
	}
	return DefaultPrefix + key
}
