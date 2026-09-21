package config

import (
	"log"
	"log/slog"
	"os"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPswd string
	DBName string
}

func GetConfig(logger *slog.Logger) Config {
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
		dbHost,
		dbPort,
		dbUser,
		dbPswd,
		dbName,
	}
}
