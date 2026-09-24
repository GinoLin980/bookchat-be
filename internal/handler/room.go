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

type RoomHandler interface {
	GetRoom(c *echo.Context) error
	GetRooms(c *echo.Context) error
	CreateRoom(c *echo.Context) error
	UpdateRoom(c *echo.Context) error
}

type roomHandler struct {
	service service.RoomService
	logger  *slog.Logger
}

func NewRoomHander(service service.RoomService, logger *slog.Logger) RoomHandler {
	return &roomHandler{
		service: service,
		logger:  logger,
	}
}

func (h *roomHandler) GetRoom(c *echo.Context) error {
	roomID := echo.PathParam[uint](c, "id")

	room, err := h.service.GetRoom(c.Request().Context(), roomID)
	if err != nil {
		if errors.Is(err, internalerror.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON()
}

func (h *roomHandler) GetRooms(c *echo.Context) error {

}

func (h *roomHandler) CreateRoom(c *echo.Context) error {

}

func (h *roomHandler) UpdateRoom(c *echo.Context) error {

}
