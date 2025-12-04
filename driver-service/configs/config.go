package configs

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database MongoConfig
}

type ServerConfig struct {
	Address string
}

type MongoConfig struct {
	URI      string
	Database string
}

func NewConfig() *Config {
	err := godotenv.Load("configs/dev.env")
	if err != nil {
		// Don't panic here, just log or ignore if using system envs
		fmt.Println("Warning: Error loading .env file")
	}

	mongoURI := buildMongoURI()

	return &Config{
		Server: ServerConfig{
			Address: getEnvOrDefault("SERVER_ADDRESS", ":8080"),
		},
		Database: MongoConfig{
			URI:      mongoURI,
			Database: getEnvOrDefault("MONGO_DB_NAME", "driverdb"),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func buildMongoURI() string {
	host := getEnvOrDefault("MONGO_HOST", "localhost")
	port := getEnvOrDefault("MONGO_PORT", "27017")
	user := getEnvOrDefault("MONGO_USER", "root")
	pass := getEnvOrDefault("MONGO_PASS", "password")

	if user != "" && pass != "" {
		return fmt.Sprintf(
			"mongodb://%s:%s@%s:%s",
			url.QueryEscape(user),
			url.QueryEscape(pass),
			host,
			port,
		)
	}

	return fmt.Sprintf("mongodb://%s:%s", host, port)
}
