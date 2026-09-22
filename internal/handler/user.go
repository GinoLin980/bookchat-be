package handler

import (
	"bookchat/internal/dto"
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/service"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

type UserHandler interface {
	Login(c *echo.Context) error
	Register(c *echo.Context) error
}

type userHandler struct {
	userService service.UserService
	logger      *slog.Logger
}

func NewUserHandler(userService service.UserService, logger *slog.Logger) UserHandler {
	return &userHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *userHandler) Login(c *echo.Context) error {
	req, err := BindAndValidate[dto.UsersRequest](c)
	if err != nil {
		return err
	}

	result, err := h.userService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, internalerror.ErrPasswordInvalid) {
			return c.JSON(http.StatusUnauthorized, nil)
		}
		if errors.Is(err, internalerror.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, nil)
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *userHandler) Register(c *echo.Context) error {
	req, err := BindAndValidate[dto.UsersRequest](c)
	if err != nil {
		return err
	}

	result, err := h.userService.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, internalerror.ErrDuplicatedError) {
			return c.JSON(http.StatusConflict, map[string]string{"message": "you have already registered"})
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, result)
}
