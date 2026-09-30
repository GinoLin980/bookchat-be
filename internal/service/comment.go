package service

import (
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/model"
	"bookchat/internal/repo"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"
)

type CommentService interface {
	CreateComment(ctx context.Context, userID, roomID uint, comment *model.Comment) error
	GetComments(ctx context.Context, roomID uint) ([]model.Comment, error)
}

type commentService struct {
	repo     repo.CommentRepo
	roomRepo repo.RoomRepo
	logger   *slog.Logger
}

func NewCommentService(repo repo.CommentRepo, roomRepo repo.RoomRepo, logger *slog.Logger) CommentService {
	return &commentService{
		repo:     repo,
		roomRepo: roomRepo,
		logger:   logger,
	}
}

func (s *commentService) CreateComment(ctx context.Context, userID, roomID uint, comment *model.Comment) error {
	room, err := s.roomRepo.GetRoomWithoutUserInfo(ctx, roomID)
	if err != nil {
		return err
	}

	// either should be moderator or (registered and being assigned and the scheduled time started)
	switch {
	// if it's moderator, then bypass
	case room.UserID == userID:
	// if it's not user's turn
	case room.AssignedToComment != userID:
		return fmt.Errorf("it's not your turn! %w", internalerror.ErrUserForbidden)

	// if user is not registered(approved by moderator)
	case !slices.Contains(room.Registered, userID):
		return fmt.Errorf("You're not registered! %w", internalerror.ErrUserForbidden)

	// if the time is before starting time
	case time.Now().Before(room.ScheduledDate):
		return fmt.Errorf("the discussion is not started yet! %w", internalerror.ErrUserForbidden)

	// if the discussion has ended already
	case room.State == "ended":
		return fmt.Errorf("the discussion has ended! %w", internalerror.ErrUserForbidden)

	}

	if err := s.repo.CreateComment(ctx, userID, roomID, comment); err != nil {
		return err
	}

	return nil
}

func (s *commentService) GetComments(ctx context.Context, roomID uint) ([]model.Comment, error) {
	comments, err := s.repo.GetComments(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
