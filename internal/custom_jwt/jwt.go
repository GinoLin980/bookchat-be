package customjwt

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type CustomJWTClaims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

func NewCustomJWTClaims(username string, userID uint, expires time.Duration) *CustomJWTClaims {
	return &CustomJWTClaims{
		username,
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * expires)),
		},
	}
}

func GetClaimsFromCtx(c *echo.Context) (*CustomJWTClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid JWT token")
	}

	claims, ok := token.Claims.(*CustomJWTClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid claims type")
	}

	return claims, nil
}
