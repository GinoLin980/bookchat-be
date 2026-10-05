package database

import (
	"bookchat/internal/config"
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// form DSN and use GORM Postgres driver to connnect
func ConnectDB(config config.Config, logger *slog.Logger) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=prefer",
		config.DBHost,
		config.DBUser,
		config.DBPswd,
		config.DBName,
		config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	return db
}
