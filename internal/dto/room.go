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
	AssignedToComment uint      `json:"assigned_to_comment"`
}

type RoomReponse struct {
	RoomID     uint
	BookTitle  string    `json:"book_title"`
	BookAuthor string    `json:"book_author"`
	Moderator  RoomUser  `json:"moderator"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`

	ScheduledDate     time.Time `json:"scheduled_date"`
	AssignedToComment uint      `json:"assigned_to_comment"`

	Comments []CommentReponse `json:"comments"`
}

type RoomUser struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
}
