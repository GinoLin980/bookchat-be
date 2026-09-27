package handler

import (
	customjwt "bookchat/internal/custom_jwt"
	"bookchat/internal/dto"
	internalerror "bookchat/internal/internal_error"
	"bookchat/internal/model"
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
	roomID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return err
	}

	var userID uint = 0
	claims, ok := customjwt.TryGetClaimsFromCtx(c)
	if ok {
		userID = claims.UserID
	}

	room, err := h.service.GetRoom(c.Request().Context(), userID, roomID)
	if err != nil {
		if errors.Is(err, internalerror.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, room)
}

func (h *roomHandler) GetRooms(c *echo.Context) error {
	roomID, err := echo.PathParam[uint](c, "id")
	roomTitle := c.QueryParam("title")

	rooms, err := h.service.GetRooms(c.Request().Context(), roomID, roomTitle)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := []dto.RoomPreviewResponse{}
	for _, room := range rooms {
		resp = append(resp, room.ToPreviewResponse())
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *roomHandler) CreateRoom(c *echo.Context) error {
	claims, err := customjwt.GetClaimsFromCtx(c)
	if err != nil {
		return err
	}

	req, err := BindAndValidate[dto.RoomRequest](c)
	if err != nil {
		return err
	}

	room, err := h.service.CreateRoom(c.Request().Context(), claims.UserID, req)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, room.ToResponse(claims.UserID, []model.User{}))
}

func (h *roomHandler) UpdateRoom(c *echo.Context) error {
	claims, err := customjwt.GetClaimsFromCtx(c)
	if err != nil {
		return err
	}

	req, err := BindAndValidate[dto.RoomUpdateRequest](c)
	if err != nil {
		return err
	}

	if err := h.service.UpdateRoom(c.Request().Context(), claims.UserID, req); err != nil {
		if errors.Is(err, internalerror.ErrUserForbidden) {
			c.JSON(http.StatusForbidden, map[string]string{"message": "you're not moderator"})
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
