package customjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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
