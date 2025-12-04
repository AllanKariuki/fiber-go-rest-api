package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID 	  uint      `json:"id" gorm:"primaryKey"`
	Name   string    `json:"name"`
	Email string    `json:"email" gorm:"unique; not null"`
	Password string  `json:"-" gorm:"not null"` // Exclude password from JSON responses
	Role	string    `json:"role" gorm:"default:'user'"` // e.g., "admin" or "user"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// Hashpassword hashes the user password before saving to the database
func (u *User) Hashpassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

// ConfirmPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// DTOs
type RegisterDTO struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserDTO struct {
	Name  string `json:"name" validate:"omitempty,min=2"`
	Email string `json:"email" validate:"omitempty,email"`
}