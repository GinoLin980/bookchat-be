package route

import (
	"bookchat/internal/handler"
	"bookchat/internal/repo"
	"bookchat/internal/service"
	"log/slog"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func LoadRoutes(e *echo.Echo, secret string, db *gorm.DB, logger *slog.Logger) {

	api := e.Group("/api/v1")
	loadUserRoutes(api, secret, db, logger)
}

func loadUserRoutes(e *echo.Group, secret string, db *gorm.DB, logger *slog.Logger) {
	jwtService := service.NewJWTService(secret, logger)
	userRepo := repo.NewUserRepo(db, logger)
	userService := service.NewUserService(jwtService, userRepo, logger)
	userHandler := handler.NewUserHandler(userService, logger)

	e.POST("/login", userHandler.Login)
	e.POST("/register", userHandler.Register)
}
