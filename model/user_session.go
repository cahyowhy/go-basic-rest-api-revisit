package model

import (
	"time"
)

type UserSession struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	Expired      time.Time `json:"expired"`
	RefreshToken string    `json:"refresh_token"`
	Username     string    `gorm:"type:varchar(50);unique" json:"username"`
}
