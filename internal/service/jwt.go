package service

import (
	"log/slog"
	"time"
)

type JWTService interface {
	Issue(username string, userID int, expires time.Duration) (string, error)
}

type jwtService struct {
	logger *slog.Logger
}
