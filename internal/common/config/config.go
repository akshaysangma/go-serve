// Package config contains the Configurations need by App and
// function to load from environment variable or file.
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig       `mapstructure:"APP"`
	Log       LogConfig       `mapstructure:"LOG"`
	Database  DatabaseConfig  `mapstructure:"DATABASE"`
	JWT       JWTConfig       `mapstructure:"JWT"`
	RateLimit RateLimitConfig `mapstructure:"RATE_LIMIT"`
}

type AppConfig struct {
	Port                   int           `mapstructure:"PORT"`
	GracefulShutdownPeriod time.Duration `mapstructure:"GRACEFUL_SHUTDOWN_PERIOD"`
}

type LogConfig struct {
	Level    string `mapstructure:"LEVEL"`
	Encoding string `mapstructure:"ENCODING"`
}

type DatabaseConfig struct {
	URL            string `mapstructure:"URL"`
	MaxConnections int    `mapstructure:"MAX_CONNECTIONS"`
}

type JWTConfig struct {
	Secret             string        `mapstructure:"SECERT"`
	ExpirationDuration time.Duration `mapstructure:"EXPIRATION_DURATION"`
}

type RateLimitConfig struct {
	LimitInterval time.Duration `mapstructure:"LIMIT_INTERVAL"`
	Burst         int           `mapstructure:"BURST"`
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("GOSERVE")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Fprintf(os.Stdout, "Config file not found, using defaults and environment variables.\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error reading config file: %s\n", err)
			os.Exit(1)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to decode into struct, %v\n", err)
		os.Exit(1)
	}

	return &config
}
