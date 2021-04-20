package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"
type UserBook struct {
	gorm.Model
	BorrowDate time.Time    `json:"borrow_date"`
	ReturnDate sql.NullTime `json:"return_date,omitempty"`
	UserID     uint         `json:"user_id"`
	BookID     uint         `json:"book_id"`
	Book       Book         `json:"book,omitempty"`
	User       User         `json:"user,omitempty"`
}

type bookUserBorrows struct {
	BookId uint `json:"book_id" validate:"required,number"`
}
type UserBorrowBook struct {
	Books []bookUserBorrows `json:"books" validate:"required,min=1,max=4"`
}
