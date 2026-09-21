package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPswd string
	DBName string
}

func GetConfig() Config {
	dbHost := os.Getenv("PG_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbUser := os.Getenv("PG_USER")
	dbPswd := os.Getenv("PG_PSWD")
	dbName := os.Getenv("PG_DB")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPswd == "" || dbName == "" {
		log.Fatal("Database env var not set")
	}

	return Config{
		dbHost,
		dbPort,
		dbUser,
		dbPswd,
		dbName,
	}
}
