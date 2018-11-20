package model

import "time"

type User struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Name     string `gorm:"unique;not null" json:"name"`
	Password string `json:"-"`
}

type Chat struct {
	ID         int  `gorm:"primary_key" json:"id"`
	RoomUser   User `json:"room_user"`
	RoomUserID int
	User       User `json:"user"`
	UserID     int
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
}
