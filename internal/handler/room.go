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
	ApproveUser(c *echo.Context) error
	DenyRequest(c *echo.Context) error
	PassTurn(c *echo.Context) error
	ApplyRequest(c *echo.Context) error
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

// GetRoom Get the detail of a room
// @Summary Get the detail of a room
// @Tags rooms
// @Produce json
// @Security BearerAuth || {}
// @Param id path uint true "Room ID"
// @Success 200 {object} dto.RoomResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /rooms/{id} [get]
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

// GetRooms Get the preview of rooms or find room via room ID or keyword in title
// @Summary Get the preview of rooms or find room via room ID or keyword in title
// @Tags rooms
// @Param id query uint false "Room ID"
// @Param title query string false "Room title"
// @Produce json
// @Success 200 {object} []dto.RoomPreviewResponse
// @Failure 500
// @Router /rooms [get]
func (h *roomHandler) GetRooms(c *echo.Context) error {
	roomID, err := echo.QueryParam[uint](c, "id")
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

// CreateRoom Create a room
// @Summary Create a room
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.RoomRequest true "Create room request"
// @Success 201 {object} dto.RoomResponse
// @Failure 500
// @Router /rooms [post]
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

	// insert room.Moderator in the second param for unneccessary user find query
	return c.JSON(http.StatusCreated, room.ToResponse(claims.UserID, []model.User{room.Moderator}, []model.User{}))
}

// UpdateRoom Update the room's information or manual adding reqeust users, only moderator is allowed to this
// @Summary Update the room's information or manual adding reqeust users, only moderator is allowed to this
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "Room ID"
// @Param body body dto.RoomUpdateRequest true "Update room request"
// @Success 204
// @Failure 403 {object} map[string]string
// @Failure 500
// @Router /rooms/{id} [patch]
func (h *roomHandler) UpdateRoom(c *echo.Context) error {
	claims, err := customjwt.GetClaimsFromCtx(c)
	if err != nil {
		return err
	}

	roomID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return err
	}

	req, err := BindAndValidate[dto.RoomUpdateRequest](c)
	if err != nil {
		return err
	}

	if err := h.service.UpdateRoom(c.Request().Context(), claims.UserID, roomID, req); err != nil {
		if errors.Is(err, internalerror.ErrUserForbidden) {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "you're not moderator"})
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

// ApproveUser Move a user from requested list to registered, if user is not in requested then return 422, only moderator is allowed to this
// @Summary Move a user from requested list to registered, if user is not in requested then return 422, only moderator is allowed to this
// @Tags rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "Room ID"
// @Param body body dto.RoomApproveRequest true "Room approval request"
// @Success 204
// @Failure 422 {object} map[string]string
// @Failure 500
// @Router /rooms/{id}/approve [post]
func (h *roomHandler) ApproveUser(c *echo.Context) error {
	claims, err := customjwt.GetClaimsFromCtx(c)
	if err != nil {
		return err
	}

	roomID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	req, err := BindAndValidate[dto.RoomApproveRequest](c)
	if err != nil {
		return err
	}

	if err := h.service.ApproveRequested(c.Request().Context(), claims.UserID, req.ApproveUserID, roomID); err != nil {
		if errors.Is(err, internalerror.ErrUserForbidden) {
			return c.NoContent(http.StatusForbidden)
		} else if errors.Is(err, internalerror.ErrUnprocessableEntity) {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"message": "user not in request list, ask user to request the room first"})
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

// DenyRequest Deny a user's request to join room, only available for moderator
// @Summary Deny a user's request to join room, only available for moderator
// @Tags rooms
// @Produce json
// @Param id path uint true "Room ID"
// @Security BearerAuth
// @Success 204
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /rooms/{id}/deny [post]
func (h *roomHandler) DenyRequest(c *echo.Context) error {
	claims, err := customjwt.GetClaimsFromCtx(c)
	if err != nil {
		return err
	}

	roomID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return err
	}

	req, err := BindAndValidate[dto.RoomDenyRequest](c)
	if err != nil {
		return err
	}

	if err := h.service.DenyRequest(c.Request().Context(), claims.UserID, req.DenyUserID, roomID); err != nil {
		if errors.Is(err, internalerror.ErrUserForbidden) {
			return c.JSON(http.StatusForbidden, err)
		}
		if errors.Is(err, internalerror.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

// PassTurn Pass the current user's turn if it's current assigned, if not, return 404, or it might be the room not found
// @Summary  Pass the current user's turn if it's current assigned, if not, return 404, or it might be the room not found
// @Tags rooms
// @Param id path uint true "Room ID"
// @Security BearerAuth
// @Success 204
// @Failure 401
// @Failure 403
// @Failure 404
// @Failure 500
// @Router /rooms/{id}/pass [post]
func (h *roomHandler) PassTurn(c *echo.Context) error {
	claims, err := customjwt.GetClaimsFromCtx(c)
	if err != nil {
		return err
	}

	roomID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return err
	}

	if err := h.service.PassTurn(c.Request().Context(), claims.UserID, roomID); err != nil {
		if errors.Is(err, internalerror.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

// ApplyRequest Apply for a room's registry, if already applied, still return 200
// @Summary Apply for a room's registry, if already applied, still return 200
// @Tags rooms
// @Produce json
// @Param id path uint true "Room ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Failure 404
// @Failure 500
// @Router /rooms/{id}/apply [post]
func (h *roomHandler) ApplyRequest(c *echo.Context) error {
	claims, err := customjwt.GetClaimsFromCtx(c)
	if err != nil {
		return err
	}

	roomID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.service.ApplyRequested(c.Request().Context(), claims.UserID, roomID); err != nil {
		if errors.Is(err, internalerror.ErrRecordNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
