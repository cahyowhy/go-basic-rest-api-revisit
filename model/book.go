package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title        string         `gorm:"type:varchar(60)" json:"title"`
	Sheet        uint           `json:"sheet"`
	DateOffIssue time.Time      `json:"date_off_issue,omitempty"`
	Introduction string         `json:"introduction,omitempty"`
	Author       pq.StringArray `gorm:"type:varchar[]" json:"author"`
}
