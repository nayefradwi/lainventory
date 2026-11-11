package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Environment string

const (
	EnvironmentProduction  Environment = "production"
	EnvironmentStaging     Environment = "staging"
	EnvironmentDevelopment Environment = "development"
)

type Config struct {
	Environment Environment
	Address     string
	DatabaseUrl string
	RedisUrl    string
}

func NewConfig() Config {
	godotenv.Load()
	address := os.Getenv("ADDRESS")
	if address == "" {
		address = ":4000"
	}

	return Config{
		Environment: Environment(os.Getenv("ENVIRONMENT")),
		Address:     address,
		DatabaseUrl: os.Getenv("GOOSE_DBSTRING"),
		RedisUrl:    os.Getenv("REDIS_URL"),
	}
}
