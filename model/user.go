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
	FirstName   string    `gorm:"type:varchar(50)" json:"first_name" validate:"required"`
	LastName    string    `gorm:"type:varchar(50)" json:"last_name" validate:"required"`
	Username    string    `gorm:"type:varchar(50);not null;unique" json:"username" validate:"required"`
	Email       string    `gorm:"type:varchar(50);not null;unique" json:"email" validate:"required,email"`
	Password    string    `gorm:"not null" json:"password" validate:"required"`
	PhoneNumber string    `gorm:"type:varchar(50)" json:"phone_number" validate:"required"`
	BirthDate   time.Time `json:"birth_date,omitempty"`

	// gorm:"-"
	// gorm:"ForeignKey:UserID"
	UserBooks []UserBook `json:"user_books,omitempty"`
	// Run this shit first : CREATE TYPE user_role AS ENUM ( 'ADMIN', 'USER');
	// https://github.com/go-gorm/gorm/issues/1978
	UserRole UserRole `json:"user_role" validate:"required"`
}

type LoginUser struct {
	Username string `json:"username" validate:"required_without=Email"`
	Email    string `json:"email" validate:"required_without=Username"`
	Password string `json:"password" validate:"required"`
}
