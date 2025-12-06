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
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: Error loading .env file")
	}

	mongoURI := buildMongoURI()

	return &Config{
		Server: ServerConfig{
			Address: os.Getenv("SERVER_ADDRESS"),
		},
		Database: MongoConfig{
			URI:      mongoURI,
			Database: os.Getenv("MONGO_DB_NAME"),
		},
	}
}

func buildMongoURI() string {
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	user := os.Getenv("MONGO_USER")
	pass := os.Getenv("MONGO_PASS")

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
