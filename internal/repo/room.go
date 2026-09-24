package repo

import (
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/model"
	"bookchat/internal/query"
	"context"
	"errors"
	"log/slog"
	"strings"

	"gorm.io/gorm"
)

type RoomRepo interface {
	GetRoom(ctx context.Context, roomID uint) (model.Room, error)
	GetRooms(ctx context.Context, roomID uint, roomTitle string) ([]model.Room, error)
	CreateRoom(ctx context.Context, room *model.Room) error
	UpdateRoom(ctx context.Context, userID, roomID uint, room model.Room) error
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

func (r *roomRepo) GetRoom(ctx context.Context, roomID uint) (model.Room, error) {
	room, err := gorm.G[model.Room](r.db).Where(query.Room.ID.Eq(roomID)).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return room, internalerror.ErrRecordNotFound
		}
		r.logger.Error(err.Error())
		return room, internalerror.ErrDatabaseErr
	}

	return room, nil
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

func (r *roomRepo) CreateRoom(ctx context.Context, room *model.Room) error {
	if err := gorm.G[model.Room](r.db).Create(ctx, room); err != nil {
		r.logger.Error(err.Error())
		return internalerror.ErrDatabaseErr
	}

	return nil
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
