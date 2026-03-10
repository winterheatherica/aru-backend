package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                     uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email                  string     `gorm:"column:email;not null;unique" json:"email"`
	Username               *string    `gorm:"column:username;unique" json:"username,omitempty"`
	FullName               *string    `gorm:"column:full_name" json:"full_name,omitempty"`
	AvatarURL              *string    `gorm:"column:avatar_url" json:"avatar_url,omitempty"`
	Password               string     `gorm:"column:password;not null" json:"-"`
	Role                   string     `gorm:"column:role;type:user_role;not null;default:MEMBER" json:"role"`
	IsActive               bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	EmailVerified          bool       `gorm:"column:email_verified;not null;default:false" json:"email_verified"`
	EmailVerificationToken *string    `gorm:"column:email_verification_token" json:"email_verification_token,omitempty"`
	EmailVerifiedAt        *time.Time `gorm:"column:email_verified_at" json:"email_verified_at,omitempty"`
	ResetPasswordToken     *string    `gorm:"column:reset_password_token" json:"reset_password_token,omitempty"`
	ResetPasswordExpires   *time.Time `gorm:"column:reset_password_expires" json:"reset_password_expires,omitempty"`
	LastPasswordChange     *time.Time `gorm:"column:last_password_change" json:"last_password_change,omitempty"`
	LastLoginAt            *time.Time `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	LastLoginIP            *string    `gorm:"column:last_login_ip;type:inet" json:"last_login_ip,omitempty"`
	LastLoginAgent         *string    `gorm:"column:last_login_agent" json:"last_login_agent,omitempty"`
	LoginCount             int        `gorm:"column:login_count;not null;default:0" json:"login_count"`
	FailedAttempts         int        `gorm:"column:failed_attempts;not null;default:0" json:"failed_attempts"`
	LockedUntil            *time.Time `gorm:"column:locked_until" json:"locked_until,omitempty"`
	CreatedAt              time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
