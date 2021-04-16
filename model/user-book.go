package model

import (
	"time"

	"gorm.io/gorm"
)

type UserBook struct {
	gorm.Model
	BorrowDate time.Time `json:"borrow_date"`
	ReturnDate time.Time `json:"return_date,omitempty"`
	UserID     uint      `json:"user_id"`
	BookID     uint      `json:"book_id"`
	Book       Book      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"book,omitempty"`
	User       User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}
