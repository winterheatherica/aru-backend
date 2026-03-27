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

type UserAdminUpsertInput struct {
	Email         string `json:"email"`
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	Password      string `json:"password,omitempty"`
	Role          string `json:"role"`
	Active        *bool  `json:"active,omitempty"`
	EmailVerified *bool  `json:"email_verified,omitempty"`
}

type UserAdminItem struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	FullName      string    `json:"full_name"`
	Role          string    `json:"role"`
	Active        bool      `json:"active"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
