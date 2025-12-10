package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null;size:100"`
	Password  string    `json:"-" gorm:"not null;size:255"`
	Email     string    `json:"email" gorm:"unique;not null;size:150"`
	FullName  string    `json:"full_name" gorm:"not null;size:200"`
	Role      string    `json:"role" gorm:"default:staff;size:50"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
