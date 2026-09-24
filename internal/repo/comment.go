package repo

import (
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/model"
	"bookchat/internal/query"
	"context"
	"errors"
	"log/slog"
	"slices"

	"gorm.io/gorm"
)

type CommentRepo interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	// GetComments(ctx context.Context, roomID uint) ([]model.Comment, error)
}

type commentRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewCommentRepo(db *gorm.DB, logger *slog.Logger) CommentRepo {
	return &commentRepo{
		db:     db,
		logger: logger,
	}
}

func (r *commentRepo) CreateComment(ctx context.Context, comment *model.Comment) error {
	room, err := gorm.G[model.Room](r.db).
		Where(query.Room.ID.Eq(comment.RoomID)).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return internalerror.ErrRecordNotFound
		}
		r.logger.Error(err.Error())
		return internalerror.ErrDatabaseErr
	}

	if !slices.Contains(room.Registered, comment.UserID) {
		return internalerror.ErrUserForbidden
	}

	if err := gorm.G[model.Comment](r.db).Create(ctx, comment); err != nil {
		r.logger.Error(err.Error())
		return internalerror.ErrDatabaseErr
	}

	// TODO: send server-side event

	return nil
}

// func (r *commentRepo) GetComments(ctx context.Context, roomID uint) ([]model.Comment, error) {
// 	comments, err := gorm.G[model.Comment](r.db)
//
// }
