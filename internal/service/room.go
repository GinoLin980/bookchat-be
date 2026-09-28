package service

import (
	"bookchat/internal/dto"
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/model"
	"bookchat/internal/repo"
	"context"
	"log/slog"
	"slices"
)

type RoomService interface {
	GetRoom(ctx context.Context, userID, roomID uint) (dto.RoomReponse, error)
	GetRooms(ctx context.Context, roomID uint, roomTitle string) ([]model.Room, error)
	CreateRoom(ctx context.Context, userID uint, req *dto.RoomRequest) (*model.Room, error)
	UpdateRoom(ctx context.Context, userID uint, req *dto.RoomUpdateRequest) error
	ApproveRequested(ctx context.Context, userID, requestUserID, roomID uint) error
	ApplyRequested(ctx context.Context, userID, roomID uint) error
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

func (s *roomService) GetRoom(ctx context.Context, userID, roomID uint) (dto.RoomReponse, error) {
	var resp dto.RoomReponse

	room, registeredUsers, requestedUsers, err := s.repo.GetRoom(ctx, roomID)
	if err != nil {
		return resp, err
	}

	resp = room.ToResponse(userID, registeredUsers, requestedUsers)

	return resp, err
}

func (s *roomService) GetRooms(ctx context.Context, roomID uint, roomTitle string) ([]model.Room, error) {
	rooms, err := s.repo.GetRooms(ctx, roomID, roomTitle)
	if err != nil {
		return rooms, err
	}

	return rooms, err
}

func (s *roomService) CreateRoom(ctx context.Context, userID uint, req *dto.RoomRequest) (*model.Room, error) {
	room := &model.Room{
		UserID:        userID,
		Title:         req.Title,
		BookTitle:     req.BookTitle,
		BookAuthor:    req.BookAuthor,
		ScheduledDate: req.ScheduledDate,
		Registered:    []uint{userID},
	}

	resp, err := s.repo.CreateRoom(ctx, room)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *roomService) UpdateRoom(ctx context.Context, userID uint, req *dto.RoomUpdateRequest) error {
	ogRoom, _, _, err := s.repo.GetRoom(ctx, req.RoomID)
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

	// approve user into registered
	if req.ApproveUserID != 0 {
		idx := slices.Index(ogRoom.Requested, req.ApproveUserID)
		if idx == -1 { // safe guard
			return internalerror.ErrUnprocessableEntity
		}
		room.Registered = append(ogRoom.Registered, req.ApproveUserID)
		room.Requested = slices.Delete(ogRoom.Requested, idx, idx+1)
	}

	if err := s.repo.UpdateRoom(ctx, userID, req.RoomID, room); err != nil {
		return err
	}

	return nil
}

func (s *roomService) ApproveRequested(ctx context.Context, userID, requestUserID, roomID uint) error {
	if err := s.UpdateRoom(ctx, userID, &dto.RoomUpdateRequest{RoomID: roomID, ApproveUserID: requestUserID}); err != nil {
		return err
	}

	return nil
}

func (s *roomService) ApplyRequested(ctx context.Context, userID, roomID uint) error {
	return s.repo.ApplyRequest(ctx, userID, roomID)
}
