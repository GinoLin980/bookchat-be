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

// Login allow user to login
// @Summary Login via POST with JSON of username,password
// @Tags users
// @Accept json
// @Produce json
// @param body body dto.UsersRequest true "User login data"
// @Success 200 {object} dto.UsersResponse
// @Failure 400
// @Failure 404
// @Failure 401
// @Failure 500
// @Router /login [post]
func (h *userHandler) Login(c *echo.Context) error {
	req, err := BindAndValidate[dto.UsersRequest](c)
	if err != nil {
		return err
	}

	result, err := h.userService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, internalerror.ErrPasswordInvalid) {
			return c.NoContent(http.StatusUnauthorized)
		}
		if errors.Is(err, internalerror.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, result)
}

// Register allow user to register
// @Summary Register via POST with JSON of username,password
// @Tags users
// @Accept json
// @Produce json
// @param body body dto.UsersRequest true "User register data"
// @Success 201 {object} dto.UsersResponse
// @Failure 400
// @Failure 409
// @Failure 500
// @Router /register [post]
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
