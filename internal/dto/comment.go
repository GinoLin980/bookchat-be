package dto

type CommentReponse struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	Content  string `json:"content"`
}
