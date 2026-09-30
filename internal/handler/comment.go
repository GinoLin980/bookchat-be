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

type CommentHandler interface {
	GetComments(c *echo.Context) error
	CreateComment(c *echo.Context) error
}

type commentHandler struct {
	service service.CommentService
	logger  *slog.Logger
}

func NewCommentHandler(service service.CommentService, logger *slog.Logger) CommentHandler {
	return &commentHandler{
		service: service,
		logger:  logger,
	}
}

// GetComments Get the comments
// @Summary Get the comments
// @Tags comments
// @Param id path uint true "Room ID"
// @Produce json
// @Success 200 {object} []dto.CommentReponse
// @Failure 400
// @Failure 500
// @Router /rooms/{id}/comments [get]
func (h *commentHandler) GetComments(c *echo.Context) error {
	roomID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	comments, err := h.service.GetComments(c.Request().Context(), roomID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	commentsResp := []dto.CommentReponse{}
	for _, comment := range comments {
		commentsResp = append(commentsResp, comment.ToResponse())
	}

	return c.JSON(http.StatusOK, commentsResp)

}

// CreateComment Create a comment
// @Summary Create a comment
// @Description Moderator (room owner) bypasses all checks.
// @Description Other users get 403 when:
// @Description - it's not your turn
// @Description - you're not registered (not approved by moderator)
// @Description - the discussion is not started yet
// @Description - the discussion has ended
// @Tags comments
// @Accept json
// @Param id path uint true "Room ID"
// @Param body body dto.CommentRequest true "Comment"
// @Produce json
// @Success 201
// @Failure 400
// @Failure 401
// @Failure 403 {string} string "forbidden: not your turn / not registered / discussion not started yet / discussion has ended"
// @Failure 500
// @Router /rooms/{id}/comments [post]
func (h *commentHandler) CreateComment(c *echo.Context) error {
	roomID, err := echo.PathParam[uint](c, "id")
	if err != nil {
		return err
	}

	claims, err := customjwt.GetClaimsFromCtx(c)
	if err != nil {
		return err
	}

	req, err := BindAndValidate[dto.CommentRequest](c)
	if err != nil {
		return err
	}

	comment := model.Comment{
		UserID:  claims.UserID,
		RoomID:  roomID,
		Content: req.Content,
	}

	if err := h.service.CreateComment(c.Request().Context(), claims.UserID, roomID, &comment); err != nil {
		if errors.Is(err, internalerror.ErrUserForbidden) {
			return c.JSON(http.StatusForbidden, err.Error())
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}
