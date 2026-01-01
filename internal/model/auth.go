package model

import "github.com/google/uuid"

type AuthUser struct {
	ID       uuid.UUID
	Email    string
	Username string
	Role     string
}
