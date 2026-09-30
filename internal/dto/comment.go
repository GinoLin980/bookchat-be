package dto

type CommentRequest struct {
	Content string `json:"content" validate:"required"`
}

type CommentReponse struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	Content  string `json:"content"`
}
