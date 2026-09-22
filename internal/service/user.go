package service

import (
	"bookchat/internal/dto"
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/model"
	"bookchat/internal/repo"
	"context"
	"log/slog"
	"time"

	"github.com/alexedwards/argon2id"
)

type UserService interface {
	Login(ctx context.Context, req *dto.UsersRequest) (*dto.UsersResponse, error)
	Register(ctx context.Context, req *dto.UsersRequest) (*dto.UsersResponse, error)
}

type userService struct {
	jwtService JWTService
	userRepo   repo.UserRepo
	logger     *slog.Logger
}

func NewUserService(jwtService JWTService, userRepo repo.UserRepo, logger *slog.Logger) UserService {
	return &userService{
		jwtService: jwtService,
		userRepo:   userRepo,
		logger:     logger,
	}
}

var params = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16, // bytes
	KeyLength:   32, // bytes
}

func (s *userService) Login(ctx context.Context, req *dto.UsersRequest) (*dto.UsersResponse, error) {
	user, err := s.userRepo.GetUserLogin(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, user.Password)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	if !match {
		return nil, internalerror.ErrPasswordInvalid
	}

	token, err := s.jwtService.Issue(user.UserName, user.ID, time.Hour*time.Duration(24))
	if err != nil {
		return nil, internalerror.ErrJWTTokenGenerateError
	}

	result := &dto.UsersResponse{
		Status:   "success",
		Username: user.UserName,
		Token:    token,
	}

	return result, nil
}

func (s *userService) Register(ctx context.Context, req *dto.UsersRequest) (*dto.UsersResponse, error) {
	hashed, err := argon2id.CreateHash(req.Password, params)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, internalerror.ErrUserHashError
	}

	user := model.User{UserName: req.Username, Password: hashed}
	if err := s.userRepo.CreateUser(ctx, &user); err != nil {
		return nil, err
	}

	token, err := s.jwtService.Issue(user.UserName, user.ID, time.Hour*time.Duration(24))
	if err != nil {
		return nil, internalerror.ErrJWTTokenGenerateError
	}

	result := &dto.UsersResponse{
		Status:   "success",
		Username: user.UserName,
		Token:    token,
	}

	return result, nil
}
