package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config contains app-level settings.
type Config struct {
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
	JwtSecretKey         string        `mapstructure:"jwt_secret_key"`
	Port                 int           `mapstructure:"port"`
}

// PostgresConfig contains database settings.
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"db_name"`
	Timezone string `mapstructure:"timezone"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

var (
	config   Config
	postgres PostgresConfig
)

// GetConfig returns current app config.
func GetConfig() Config { return config }

// Postgres returns current postgres config.
func Postgres() PostgresConfig { return postgres }

// Load reads config from viper and populates package state.
// It keeps behavior identical while improving organization and defaults.
func Load() error {
	setDefaults()

	if err := unmarshalKey("app", &config); err != nil {
		return err
	}
	if err := unmarshalKey("postgres", &postgres); err != nil {
		return err
	}
	return nil
}

func unmarshalKey(key string, dst interface{}) error {
	return viper.UnmarshalKey(key, dst)
}

func setDefaults() {
	// Non-breaking defaults used only when values are absent
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.access_token_duration", time.Duration(900))
	viper.SetDefault("app.refresh_token_duration", time.Duration(432000))
	viper.SetDefault("postgres.ssl_mode", "disable")
}
