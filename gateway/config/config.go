package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Route struct {
	Path string `mapstructure:"path"`
}

type Service struct {
	BaseURL string  `mapstructure:"base_url"`
	Routes  []Route `mapstructure:"routes"`
}

type Config struct {
	RateLimitWindow    time.Duration      `mapstructure:"rate_limit_window"`
	RateLimitCount     int                `mapstructure:"rate_limit_count"`
	Services           map[string]Service `mapstructure:"services"`
	JWTSecretKey       string             `mapstructure:"JWT_SECRET_KEY"`
	APIKey             string             `mapstructure:"API_KEY"`
	ServicePort        string             `mapstructure:"SERVICE_PORT"`
	Env                string             `mapstructure:"ENV"`
	ElasticsearchURL   string             `mapstructure:"ELASTICSEARCH_URL"`
	ElasticsearchIndex string             `mapstructure:"ELASTICSEARCH_INDEX"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.MergeInConfig()

	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	err := viper.MergeInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	viper.AutomaticEnv()

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	return &config, nil
}
