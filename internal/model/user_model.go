package model

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username,omitempty"`
	FullName  string    `json:"full_name,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	User         *UserResponse `json:"user"`
}

// type VerifyUserRequest struct {
// 	Token string `validate:"required,max=100"`
// }

// type RegisterUserRequest struct {
// 	ID       string `json:"id" validate:"required,max=100"`
// 	Password string `json:"password" validate:"required,max=100"`
// 	Name     string `json:"name" validate:"required,max=100"`
// }

// type UpdateUserRequest struct {
// 	ID       string `json:"-" validate:"required,max=100"`
// 	Password string `json:"password,omitempty" validate:"max=100"`
// 	Name     string `json:"name,omitempty" validate:"max=100"`
// }

// type LoginUserRequest struct {
// 	ID       string `json:"id" validate:"required,max=100"`
// 	Password string `json:"password" validate:"required,max=100"`
// }

// type LogoutUserRequest struct {
// 	ID string `json:"id" validate:"required,max=100"`
// }

// type GetUserRequest struct {
// 	ID string `json:"id" validate:"required,max=100"`
// }
