package model

import (
	"bookchat/internal/dto"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"uniqueIndex"`
	Password string `gorm:"not null"`
}

func (u User) ToRoomUserReponse() dto.RoomUser {
	return dto.RoomUser{
		Username: u.UserName,
		UserID:   u.ID,
	}
}
