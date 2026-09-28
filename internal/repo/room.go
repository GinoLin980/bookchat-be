package repo

import (
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/model"
	"bookchat/internal/query"
	"context"
	"errors"
	"log/slog"
	"slices"
	"strings"

	"gorm.io/gorm"
)

type RoomRepo interface {
	GetRoom(ctx context.Context, roomID uint) (model.Room, []model.User, []model.User, error)
	GetRooms(ctx context.Context, roomID uint, roomTitle string) ([]model.Room, error)
	CreateRoom(ctx context.Context, room *model.Room) (model.Room, error)
	UpdateRoom(ctx context.Context, userID, roomID uint, room model.Room) error
	ApplyRequest(ctx context.Context, userID, roomID uint) error
}

type roomRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewRoomRepo(db *gorm.DB, logger *slog.Logger) RoomRepo {
	return &roomRepo{
		db:     db,
		logger: logger,
	}
}

func (r *roomRepo) GetRoom(ctx context.Context, roomID uint) (model.Room, []model.User, []model.User, error) {
	room, err := gorm.G[model.Room](r.db).
		Preload("Moderator", nil).
		Preload("Comments", nil).
		Where(query.Room.ID.Eq(roomID)).
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return room, nil, nil, internalerror.ErrRecordNotFound
		}
		r.logger.Error(err.Error())
		return room, nil, nil, internalerror.ErrDatabaseErr
	}

	registeredUsers, err := gorm.G[model.User](r.db).Where(query.User.ID.In(room.Registered...)).Find(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return room, nil, nil, internalerror.ErrDatabaseErr
	}

	requestedUsers, err := gorm.G[model.User](r.db).Where(query.User.ID.In(room.Requested...)).Find(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return room, nil, nil, internalerror.ErrDatabaseErr
	}

	return room, registeredUsers, requestedUsers, nil
}

func (r *roomRepo) GetRooms(ctx context.Context, roomID uint, roomTitle string) ([]model.Room, error) {
	roomTitle = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(roomTitle)
	if roomTitle != "" {
		roomTitle = "%" + roomTitle + "%" // fuzzy search
	} else {
		roomTitle = "%" // match all
	}

	qry := gorm.G[model.Room](r.db).Where(query.Room.Title.ILike(roomTitle))

	if roomID != 0 {
		qry = qry.Where(query.Room.ID.Eq(roomID))
	}

	rooms, err := qry.Find(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, internalerror.ErrDatabaseErr
	}

	return rooms, nil
}

func (r *roomRepo) CreateRoom(ctx context.Context, room *model.Room) (model.Room, error) {
	if err := gorm.G[model.Room](r.db).Create(ctx, room); err != nil {
		r.logger.Error(err.Error())
		return model.Room{}, internalerror.ErrDatabaseErr
	}

	resp, err := gorm.G[model.Room](r.db).
		Preload("Moderator", nil).
		Where(query.Room.ID.Eq(room.ID)).
		First(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return model.Room{}, internalerror.ErrDatabaseErr
	}

	return resp, nil
}

func (r *roomRepo) UpdateRoom(ctx context.Context, userID, roomID uint, room model.Room) error {
	rowsAffected, err := gorm.G[model.Room](r.db).
		Where(query.Room.ID.Eq(roomID)).
		Where(query.Room.UserID.Eq(userID)). // moderatorID
		Updates(ctx, room)

	if err != nil {
		r.logger.Error(err.Error())
		return internalerror.ErrDatabaseErr
	}
	if rowsAffected == 0 {
		return internalerror.ErrUserForbidden
	}

	return nil
}

func (r *roomRepo) ApplyRequest(ctx context.Context, userID, roomID uint) error {
	room, err := gorm.G[model.Room](r.db).Where(query.Room.ID.Eq(roomID)).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return internalerror.ErrRecordNotFound
		}
		r.logger.Error(err.Error())
		return internalerror.ErrDatabaseErr
	}

	if slices.Contains(room.Requested, userID) || slices.Contains(room.Registered, userID) {
		r.logger.Info("user already in requested||registered")
		return nil // no-op idempotency
	}
	requestedUsers := append(room.Requested, userID)
	updateRoom := model.Room{Requested: requestedUsers}

	if _, err := gorm.G[model.Room](r.db).Where(query.Room.ID.Eq(roomID)).Updates(ctx, updateRoom); err != nil {
		r.logger.Error(err.Error())
		return internalerror.ErrDatabaseErr
	}

	return nil
}
