package config

import (
	"github.com/joho/godotenv"
	"github.com/nayefradwi/lainventory/internal/auth"
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
	Auth        auth.AuthConfig
}

func NewConfig() Config {
	godotenv.Load()
	address := os.Getenv("ADDRESS")
	if address == "" {
		address = ":4000"
	}

	databaseUrl := os.Getenv("GOOSE_DBSTRING")
	if databaseUrl == "" {
		panic("GOOSE_DBSTRING is required")
	}

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		panic("REDIS_URL is required")
	}

	return Config{
		Environment: Environment(os.Getenv("ENVIRONMENT")),
		Address:     address,
		DatabaseUrl: databaseUrl,
		RedisUrl:    redisUrl,
		Auth:        auth.NewAuthConfig(),
	}
}
