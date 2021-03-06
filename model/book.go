package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title        string         `gorm:"type:varchar(60)" json:"title" validate:"required"`
	Sheet        uint           `json:"sheet" validate:"required,number"`
	DateOffIssue time.Time      `json:"date_off_issue,omitempty"`
	Introduction string         `json:"introduction,omitempty"`
	Author       pq.StringArray `gorm:"type:varchar[]" json:"author" validate:"required"`
}
