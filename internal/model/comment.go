package model

import (
	"bookchat/internal/dto"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID  uint
	Poster  User `gorm:"not null;foreignKey:UserID"`
	RoomID  uint
	Content string
}

func (c Comment) ToResponse() dto.CommentReponse {
	return dto.CommentReponse{
		Username: c.Poster.UserName,
		UserID:   c.UserID,
		Content:  c.Content,
	}
}
