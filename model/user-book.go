package model

import (
	"time"

	"gorm.io/gorm"
)

// gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"
type UserBook struct {
	gorm.Model
	BorrowDate time.Time `json:"borrow_date"`
	ReturnDate time.Time `json:"return_date,omitempty"`
	UserID     uint      `json:"user_id"`
	BookID     uint      `json:"book_id"`
	Book       Book      `json:"book,omitempty"`
	User       User      `json:"user,omitempty"`
}
