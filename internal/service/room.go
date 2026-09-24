package service

import (
	"bookchat/internal/dto"
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/model"
	"bookchat/internal/repo"
	"context"
	"log/slog"
)

type RoomService interface {
	GetRoom(ctx context.Context, roomID uint) (model.Room, error)
	GetRooms(ctx context.Context, roomID uint, roomTitle string) ([]model.Room, error)
	CreateRoom(ctx context.Context, userID uint, req dto.RoomRequest) error
	UpdateRoom(ctx context.Context, userID, roomID uint, req dto.RoomUpdateRequest) error
}

type roomService struct {
	repo   repo.RoomRepo
	logger *slog.Logger
}

func NewRoomService(repo repo.RoomRepo, logger *slog.Logger) RoomService {
	return &roomService{
		repo:   repo,
		logger: logger,
	}
}

func (s *roomService) GetRoom(ctx context.Context, roomID uint) (model.Room, error) {
	room, err := s.repo.GetRoom(ctx, roomID)
	if err != nil {
		return room, err
	}

	return room, err
}

func (s *roomService) GetRooms(ctx context.Context, roomID uint, roomTitle string) ([]model.Room, error) {
	rooms, err := s.repo.GetRooms(ctx, roomID, roomTitle)
	if err != nil {
		return rooms, err
	}

	return rooms, err
}

func (s *roomService) CreateRoom(ctx context.Context, userID uint, req dto.RoomRequest) error {
	room := &model.Room{
		UserID:        userID,
		Title:         req.Title,
		BookTitle:     req.BookTitle,
		BookAuthor:    req.BookAuthor,
		ScheduledDate: req.ScheduledDate,
	}

	if err := s.repo.CreateRoom(ctx, room); err != nil {
		return err
	}

	return nil
}

func (s *roomService) UpdateRoom(ctx context.Context, userID, roomID uint, req dto.RoomUpdateRequest) error {
	ogRoom, err := s.GetRoom(ctx, roomID)
	if err != nil {
		return err
	}

	if ogRoom.UserID != userID {
		return internalerror.ErrUserForbidden
	}

	room := model.Room{
		Title:             req.Title,
		BookTitle:         req.BookTitle,
		BookAuthor:        req.BookAuthor,
		AssignedToComment: ogRoom.AssignedToComment,
	}

	if req.AddUserID != 0 {
		room.Registered = append(ogRoom.Registered, userID)
	}

	if err := s.repo.UpdateRoom(ctx, userID, roomID, room); err != nil {
		return err
	}

	return nil
}
