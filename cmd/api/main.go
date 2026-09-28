package main

import (
	"bookchat/internal/config"
	customvalidator "bookchat/internal/custom_validator"
	"bookchat/internal/database"
	"bookchat/internal/model"
	"bookchat/internal/route"
	"errors"
	"log/slog"
	"net/http"
	"os"

	_ "bookchat/docs"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title BookChat Backend API
// @version 1.0
// @host bookchat.gcp.ginol.in
// @Basepath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                "Bearer {JWT}"
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Error("failed to load .env")
		os.Exit(1)
	}

	config := config.GetConfig(logger)

	db := database.ConnectDB(config, logger)
	db.AutoMigrate(&model.User{}, &model.Room{}, &model.Comment{})

	e := echo.New()

	e.Validator = &customvalidator.CustomValidator{V: customvalidator.InitValidator()}
	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		if resp, rerr := echo.UnwrapResponse(c.Response()); rerr == nil && resp != nil && resp.Committed {
			return
		}
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string, len(ve))
			for _, fe := range ve {
				out[fe.Field()] = "failed: " + fe.Tag()
			}
			c.JSON(http.StatusBadRequest, out)
			return
		}
		echo.DefaultHTTPErrorHandler(true)(c, err)
	}
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	route.LoadRoutes(e, config.JWTSecret, db, logger)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"hello": "world"})
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("port might be used, can't start server")
	}
}
