package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserFineHistory struct {
	gorm.Model
	UserID      uint          `json:"user_id"`
	User        User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
	Fine        uint          `json:"fine"`
	HasPaid     bool          `json:"has_paid"`
	UserBookIds pq.Int32Array `gorm:"type:integer[]" json:"user_book_ids"`
}
