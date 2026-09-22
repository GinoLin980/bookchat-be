package main

import (
	"bookchat/internal/config"
	"bookchat/internal/database"
	"bookchat/internal/route"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Error("failed to load .env")
		os.Exit(1)
	}

	config := config.GetConfig(logger)

	db := database.ConnectDB(config, logger)

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())

	route.LoadRoutes(e, config.JWTSecret, db, logger)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"hello": "world"})
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("port might be used, can't start server")
	}
}
