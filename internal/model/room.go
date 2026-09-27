package model

import (
	"bookchat/internal/dto"
	"slices"
	"time"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	UserID     uint
	Moderator  User   `gorm:"not null;foreignKey:UserID"`
	Registered []uint `gorm:"serializer:json"`

	Title      string `gorm:"not null"` // room's title
	BookTitle  string `gorm:"not null"`
	BookAuthor string `gorm:"not null"`

	Comments []Comment `gorm:"foreignKey:RoomID"`

	ScheduledDate     time.Time
	State             string `gorm:"oneof:started,finished"`
	AssignedToComment uint
}

func (r Room) ToPreviewResponse() dto.RoomPreviewResponse {
	return dto.RoomPreviewResponse{
		RoomID:        r.ID,
		Title:         r.Title,
		BookTitle:     r.BookTitle,
		BookAuthor:    r.BookAuthor,
		Moderator:     r.Moderator.ToRoomUserReponse(),
		ScheduledDate: r.ScheduledDate,
	}
}

func (r Room) ToResponse(userID uint, registeredUsers []User) dto.RoomReponse {
	role := "general"
	if userID == r.UserID {
		role = "moderator"
	} else if slices.Contains(r.Registered, userID) {
		role = "registered"
	}

	comments := []dto.CommentReponse{}
	for _, comment := range r.Comments {
		comments = append(comments, comment.ToResponse())
	}

	roomRegisteredUsers := []dto.RoomUser{}
	for _, user := range registeredUsers {
		roomRegisteredUsers = append(roomRegisteredUsers, user.ToRoomUserReponse())
	}

	result := dto.RoomReponse{
		RoomID:     r.ID,
		Title:      r.Title,
		BookTitle:  r.BookTitle,
		BookAuthor: r.BookAuthor,
		Moderator:  r.Moderator.ToRoomUserReponse(),
		Role:       role,
		Registered: roomRegisteredUsers,
		CreatedAt:  r.CreatedAt,

		ScheduledDate: r.ScheduledDate,
		Comments:      comments,
	}

	return result
}
