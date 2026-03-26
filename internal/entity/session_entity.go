package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID           uuid.UUID  `gorm:"column:user_id;type:uuid;not null" json:"user_id"`
	AccessToken      *string    `gorm:"column:access_token" json:"access_token,omitempty"`
	RefreshTokenHash string     `gorm:"column:refresh_token_hash;not null" json:"-"`
	JWTID            string     `gorm:"column:jwt_id;not null" json:"jwt_id"`
	IPAddress        *string    `gorm:"column:ip_address;type:inet" json:"ip_address,omitempty"`
	UserAgent        *string    `gorm:"column:user_agent" json:"user_agent,omitempty"`
	DeviceName       *string    `gorm:"column:device_name" json:"device_name,omitempty"`
	Location         *string    `gorm:"column:location" json:"location,omitempty"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	LastUsedAt       *time.Time `gorm:"column:last_used_at" json:"last_used_at,omitempty"`
	ExpiresAt        time.Time  `gorm:"column:expires_at;not null" json:"expires_at"`
	Revoked          bool       `gorm:"column:revoked;not null;default:false" json:"revoked"`
	RevokedAt        *time.Time `gorm:"column:revoked_at" json:"revoked_at,omitempty"`
	RevokedReason    *string    `gorm:"column:revoked_reason;type:revoked_reason" json:"revoked_reason,omitempty"`
	Fingerprint      *string    `gorm:"column:fingerprint" json:"fingerprint,omitempty"`
	IsCurrent        bool       `gorm:"column:is_current;default:false" json:"is_current"`
	User             User       `gorm:"foreignKey:UserID" json:"user"`
}

func (Session) TableName() string {
	return "sessions"
}
