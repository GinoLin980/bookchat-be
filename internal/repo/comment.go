package repo

import (
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/model"
	"bookchat/internal/query"
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type CommentRepo interface {
	CreateComment(ctx context.Context, userID, roomID uint, comment *model.Comment) error
	GetComments(ctx context.Context, roomID uint) ([]model.Comment, error)
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

func (r *commentRepo) CreateComment(ctx context.Context, userID, roomID uint, comment *model.Comment) error {
	if err := gorm.G[model.Comment](r.db).Create(ctx, comment); err != nil {
		r.logger.Error(err.Error())
		return internalerror.ErrDatabaseErr
	}

	return nil
}

func (r *commentRepo) GetComments(ctx context.Context, roomID uint) ([]model.Comment, error) {
	comments, err := gorm.G[model.Comment](r.db).
		Where(query.Comment.RoomID.Eq(roomID)).
		Order(query.Comment.CreatedAt.Desc()).Find(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return nil, internalerror.ErrDatabaseErr
	}

	return comments, nil
}
