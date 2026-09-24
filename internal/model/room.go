package model

import (
	"bookchat/internal/dto"
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

func (rr Room) ToResponse() dto.RoomReponse {
	result := 
}
