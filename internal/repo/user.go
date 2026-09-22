package repo

import (
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/model"
	"bookchat/internal/query"
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type UserRepo interface {
	GetUserLogin(ctx context.Context, username string) (model.User, error)
	GetUser(ctx context.Context, userID int) (model.User, error)
	GetUsersByName(ctx context.Context, username string) ([]model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type userRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepo(db *gorm.DB, logger *slog.Logger) UserRepo {
	return &userRepo{
		db:     db,
		logger: logger,
	}
}

func (r *userRepo) GetUserLogin(ctx context.Context, username string) (model.User, error) {
	user, err := gorm.G[model.User](r.db).Where(query.User.UserName.Eq(username)).First(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, internalerror.ErrRecordNotFound
		}
		return user, internalerror.ErrDatabaseErr
	}

	return user, nil
}

func (r *userRepo) GetUser(ctx context.Context, userID int) (model.User, error) {
	user, err := gorm.G[model.User](r.db).Where(query.User.ID.Eq(uint(userID))).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error(err.Error())
			return user, internalerror.ErrRecordNotFound
		}
		return user, internalerror.ErrDatabaseErr
	}

	return user, nil
}

func (r *userRepo) GetUsersByName(ctx context.Context, username string) ([]model.User, error) {
	users, err := gorm.G[model.User](r.db).Where(query.User.UserName.ILike(username)).Find(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error(err.Error())
			return nil, internalerror.ErrRecordNotFound
		}
		return nil, internalerror.ErrDatabaseErr
	}

	return users, nil
}

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) error {
	if err := gorm.G[model.User](r.db).Create(ctx, user); err != nil {
		r.logger.Error(err.Error())
		return internalerror.ErrDatabaseErr
	}

	return nil
}
