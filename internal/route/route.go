package route

import (
	customjwt "bookchat/internal/custom_jwt"
	"bookchat/internal/handler"
	"bookchat/internal/repo"
	"bookchat/internal/service"
	"errors"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func LoadRoutes(e *echo.Echo, secret string, db *gorm.DB, logger *slog.Logger) {
	// base endpoint
	api := e.Group("/api/v1")

	// jwt dep
	jwtService := service.NewJWTService(secret, logger)

	// unprotected endpoints
	loadUserRoutes(api, jwtService, db, logger)

	// load JWT configs
	jwtConfig, optionalJWTConfig := getJWTConfigs(secret)

	// protected endpoints
	enforced, optional := getProtectedRoutes(api, jwtConfig, optionalJWTConfig)
	loadRoomRoutes(db, logger, api, enforced, optional)

	loadCommentRoutes(db, logger, enforced)
}

func loadUserRoutes(e *echo.Group, jwtService service.JWTService, db *gorm.DB, logger *slog.Logger) {
	userRepo := repo.NewUserRepo(db, logger)
	userService := service.NewUserService(jwtService, userRepo, logger)
	userHandler := handler.NewUserHandler(userService, logger)

	e.POST("/login", userHandler.Login)
	e.POST("/register", userHandler.Register)

}

func getProtectedRoutes(e *echo.Group, config, optionalConfig echojwt.Config) (*echo.Group, *echo.Group) {
	// JWT enforced endpoints
	r := e.Group("")
	r.Use(echojwt.WithConfig(config))

	// optionalJWTConfig endpoints should
	o := e.Group("")
	o.Use(echojwt.WithConfig(optionalConfig))

	return r, o
}

func loadRoomRoutes(db *gorm.DB, logger *slog.Logger, nonProtected, enforced, optional *echo.Group) {
	roomRepo := repo.NewRoomRepo(db, logger)
	roomService := service.NewRoomService(roomRepo, logger)
	roomHandler := handler.NewRoomHander(roomService, logger)

	nonProtected.GET("/rooms", roomHandler.GetRooms)

	optional.GET("/rooms/:id", roomHandler.GetRoom)

	enforced.POST("/rooms", roomHandler.CreateRoom)
	enforced.PATCH("/rooms/:id", roomHandler.UpdateRoom)

	enforced.POST("/rooms/:id/apply", roomHandler.ApplyRequest)
	enforced.POST("/rooms/:id/approve", roomHandler.ApproveUser)

	enforced.POST("/rooms/:id/deny", roomHandler.DenyRequest)
	enforced.POST("/rooms/:id/pass", roomHandler.PassTurn)
}

func loadCommentRoutes(db *gorm.DB, logger *slog.Logger, config *echo.Group) {
	commentRepo := repo.NewCommentRepo(db, logger)
	roomRepo := repo.NewRoomRepo(db, logger)
	commentService := service.NewCommentService(commentRepo, roomRepo, logger)
	commentHandler := handler.NewCommentHandler(commentService, logger)

	config.GET("/rooms/:id/comments", commentHandler.GetComments)
	config.POST("/rooms/:id/comments", commentHandler.CreateComment)
}

func getJWTConfigs(secret string) (echojwt.Config, echojwt.Config) {
	config := echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(customjwt.CustomJWTClaims)
		},
		SigningKey: []byte(secret),
	}

	optionalConfig := echojwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(customjwt.CustomJWTClaims)
		},
		SigningKey:             []byte(secret),
		ContinueOnIgnoredError: true,
		ErrorHandler: func(c *echo.Context, err error) error {
			var extractionErr *echojwt.TokenExtractionError
			if errors.As(err, &extractionErr) {
				// No Authorization header at all -> anonymous, let it through.
				return nil
			}
			// Header was present but the token itself is bad -> real 401.
			return echojwt.ErrJWTInvalid
		},
	}

	return config, optionalConfig
}
