package config

import (
	"log/slog"
	"os"
)

type Config struct {
	JWTSecret string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPswd    string
	DBName    string
}

func GetConfig(logger *slog.Logger) Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	dbHost := os.Getenv("PG_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbUser := os.Getenv("PG_USER")
	dbPswd := os.Getenv("PG_PSWD")
	dbName := os.Getenv("PG_DB")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPswd == "" || dbName == "" {
		slog.Error("Database env var not set")
		os.Exit(1)
	}

	return Config{
		jwtSecret,
		dbHost,
		dbPort,
		dbUser,
		dbPswd,
		dbName,
	}
}
