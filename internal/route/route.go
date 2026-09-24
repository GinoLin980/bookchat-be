package route

import (
	customjwt "bookchat/internal/custom_jwt"
	"bookchat/internal/handler"
	"bookchat/internal/repo"
	"bookchat/internal/service"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func LoadRoutes(e *echo.Echo, secret string, db *gorm.DB, logger *slog.Logger) {
	api := e.Group("/api/v1")

	config := echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(customjwt.CustomJWTClaims)
		},
		SigningKey: []byte(secret),
	}

	jwtService := service.NewJWTService(secret, logger)

	loadUserRoutes(api, jwtService, db, logger)
	loadProtectedRoutes(api, config)
}

func loadUserRoutes(e *echo.Group, jwtService service.JWTService, db *gorm.DB, logger *slog.Logger) {
	userRepo := repo.NewUserRepo(db, logger)
	userService := service.NewUserService(jwtService, userRepo, logger)
	userHandler := handler.NewUserHandler(userService, logger)

	e.POST("/login", userHandler.Login)
	e.POST("/register", userHandler.Register)
}

func loadProtectedRoutes(e *echo.Group, config echojwt.Config) {
	r := e.Group("")
	r.Use(echojwt.WithConfig(config))

	r.GET("/hello", func(c *echo.Context) error {
		claims, err := customjwt.GetClaimsFromCtx(c)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{"message": claims.Username})
	})
}
