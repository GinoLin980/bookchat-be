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
	// context to receive SIGKILL from host
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	// defer(stop in the end)
	defer stop()

	// initialize a logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// .env loader for local development
	if err := godotenv.Load(); err != nil {
		logger.Warn("failed to load .env")
	}

	// get the environment variables
	config := config.GetConfig(logger)

	// connect the database
	db := database.ConnectDB(config, logger)
	// db.AutoMigrate(&model.User{}, &model.Room{}, &model.Comment{})

	// Echo server
	e := echo.New()

	// custom validator by PlayGround and better 400 response format
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

	// Disable CORS
	e.Pre(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
		MaxAge:       3600,
	}))
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover()) // won't stop the server after panic
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))) // rate limiter so it won't burn my wallet
	e.Use(middleware.BodyLimit(1 << 20))

	// load all routes
	route.LoadRoutes(e, config.JWTSecret, db, logger)

	// easy endpoint for testing
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(200, map[string]string{"hello": "world"})
	})

	// API docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// get the port env, required by Cloud Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// start config including the GracefulTimeout of 5 sec
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
