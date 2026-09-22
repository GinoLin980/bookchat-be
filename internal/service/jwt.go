package service

import (
	customjwt "bookchat/internal/custom_jwt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	Issue(username string, userID uint, expires time.Duration) (string, error)
}

type jwtService struct {
	secret string
	logger *slog.Logger
}

func NewJWTService(secret string, logger *slog.Logger) JWTService {
	if secret == "" {
		logger.Error("JWT secret not provided!")
		return nil
	}
	return &jwtService{
		secret: secret,
		logger: logger,
	}
}

func (s *jwtService) Issue(username string, userID uint, expires time.Duration) (string, error) {
	claims := customjwt.NewCustomJWTClaims(username, userID, expires)

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := rawToken.SignedString([]byte(s.secret))
	if err != nil {
		s.logger.Error(err.Error())
		return "", err
	}

	return token, nil
}
