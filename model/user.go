package model

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	ADMIN UserRole = "ADMIN"
	USER  UserRole = "USER"
)

func (u *UserRole) Scan(value interface{}) error {
	*u = UserRole(value.(string))
	return nil
}

func (u UserRole) Value() (driver.Value, error) {
	return string(u), nil
}

type User struct {
	gorm.Model
	FirstName   string     `gorm:"type:varchar(20)" json:"first_name"`
	LastName    string     `gorm:"type:varchar(20)" json:"last_name"`
	Username    string     `gorm:"type:varchar(50);not null;unique" json:"username"`
	Email       string     `gorm:"type:varchar(50);not null;unique" json:"email"`
	Password    string     `gorm:"not null" json:"password"`
	PhoneNumber string     `gorm:"type:varchar(20)" json:"phone_number"`
	BirthDate   time.Time  `json:"birth_date,omitempty"`
	UserBooks   []UserBook `gorm:"-" json:"user_books,omitempty"`
	// Run this shit first : CREATE TYPE user_role AS ENUM ( 'ADMIN', 'USER');
	// https://github.com/go-gorm/gorm/issues/1978
	UserRole UserRole `gorm:"type:user_role" json:"user_role"`
}
