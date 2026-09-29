package main

import (
	"bookchat/internal/config"
	customvalidator "bookchat/internal/custom_validator"
	"bookchat/internal/database"
	"bookchat/internal/route"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Warn("failed to load .env")
	}

	config := config.GetConfig(logger)

	db := database.ConnectDB(config, logger)
	// db.AutoMigrate(&model.User{}, &model.Room{}, &model.Comment{})

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
	e.Pre(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
		MaxAge:       3600,
	}))
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.BodyLimit(1 << 20))

	route.LoadRoutes(e, config.JWTSecret, db, logger)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"hello": "world"})
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	sc := echo.StartConfig{
		Address:         ":" + port,
		HideBanner:      true,
		GracefulTimeout: 8 * time.Second,
		BeforeServeFunc: func(s *http.Server) error {
			s.ReadHeaderTimeout = 5 * time.Second
			s.IdleTimeout = 60 * time.Second
			return nil
		},
	}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("port might be used, can't start server")
	}
}
