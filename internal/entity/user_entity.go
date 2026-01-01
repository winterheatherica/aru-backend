package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email                  string    `gorm:"unique;not null"`
	Username               string    `gorm:"unique"`
	FullName               string
	AvatarURL              string
	Password               string `gorm:"not null"`
	Role                   string `gorm:"type:user_role;default:'member'"`
	IsActive               bool   `gorm:"default:true"`
	EmailVerified          bool   `gorm:"default:false"`
	EmailVerificationToken string
	EmailVerifiedAt        *time.Time
	ResetPasswordToken     string
	ResetPasswordExpires   *time.Time
	LastPasswordChange     *time.Time
	LastLoginAt            *time.Time
	LastLoginIP            string
	LastLoginAgent         string
	LoginCount             int `gorm:"default:0"`
	FailedAttempts         int `gorm:"default:0"`
	LockedUntil            *time.Time
	CreatedAt              time.Time `gorm:"autoCreateTime"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}
