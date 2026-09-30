package dto

import "time"

type RoomRequest struct {
	Title         string    `json:"title" validate:"required"`
	BookTitle     string    `json:"book_title" validate:"required"`
	BookAuthor    string    `json:"book_author" validate:"required"`
	ScheduledDate time.Time `json:"scheduled_date" validate:"required"`
}

type RoomUpdateRequest struct {
	Title             string    `json:"title"`
	BookTitle         string    `json:"book_title"`
	BookAuthor        string    `json:"book_author"`
	ScheduledDate     time.Time `json:"scheduled_date"`
	AddUserID         uint      `json:"add_user_id"`
	ApproveUserID     uint      `json:"approve_user_id"`
	AssignedToComment uint      `json:"assigned_to_comment"`
}

type RoomApproveRequest struct {
	ApproveUserID uint `json:"approve_user_id" validate:"required"`
}

type RoomPreviewResponse struct {
	RoomID        uint      `json:"room_id"`
	Title         string    `json:"title"`
	BookTitle     string    `json:"book_title"`
	BookAuthor    string    `json:"book_author"`
	Moderator     RoomUser  `json:"moderator"`
	ScheduledDate time.Time `json:"scheduled_date"`
}

type RoomResponse struct {
	RoomID     uint       `json:"room_id" validate:"required"`
	Title      string     `json:"title"`
	BookTitle  string     `json:"book_title"`
	BookAuthor string     `json:"book_author"`
	Moderator  RoomUser   `json:"moderator"`
	Registered []RoomUser `json:"registered"`
	Requested  []RoomUser `json:"requested"`
	Role       string     `json:"role"`
	CreatedAt  time.Time  `json:"created_at"`

	ScheduledDate     time.Time `json:"scheduled_date"`
	AssignedToComment uint      `json:"assigned_to_comment"`

	Comments []CommentReponse `json:"comments"`
}

type RoomUser struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
}
