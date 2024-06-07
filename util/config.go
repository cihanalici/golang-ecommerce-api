package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config is the configuration for the application
// It should be loaded from a file or other sources
type Config struct {
	DBDriver                 string        `mapstructure:"DB_DRIVER"`
	DBSource                 string        `mapstructure:"DB_SOURCE"`
	ServerAddress            string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey        string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration      time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	ResetPasswordDuration    time.Duration `mapstructure:"RESET_PASSWORD_DURATION"`
	ResetPasswordRedirectURL string        `mapstructure:"RESET_PASSWORD_REDIRECT_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
