package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID  uint
	Poster  User `gorm:"not null;foreignKey:UserID"`
	RoomID  uint
	Content string
}
