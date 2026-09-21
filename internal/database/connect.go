package database

import (
	"bookchat/internal/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(config config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=enabled",
		config.DBHost,
		config.DBUser,
		config.DBPswd,
		config.DBName,
		config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
