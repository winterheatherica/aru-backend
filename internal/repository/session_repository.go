package repository

import (
	"aru-backend/internal/entity"
	"time"

	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(session *entity.Session) error
	FindActiveByAccessToken(token string) (*entity.Session, error)
	RevokeByAccessToken(token string, reason string) error
	UpdateLastUsed(sessionID string) error
}

type sessionRepositoryImpl struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepositoryImpl{db: db}
}

func (r *sessionRepositoryImpl) Create(session *entity.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepositoryImpl) FindActiveByAccessToken(token string) (*entity.Session, error) {
	var session entity.Session
	err := r.db.Preload("User").Where("access_token = ? AND revoked = false AND expires_at > now()", token).Take(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepositoryImpl) RevokeByAccessToken(token string, reason string) error {
	now := time.Now()
	return r.db.Model(&entity.Session{}).
		Where("access_token = ? AND revoked = false", token).
		Updates(map[string]any{
			"revoked":        true,
			"revoked_at":     &now,
			"revoked_reason": reason,
			"is_current":     false,
		}).Error
}

func (r *sessionRepositoryImpl) UpdateLastUsed(sessionID string) error {
	now := time.Now()
	return r.db.Model(&entity.Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]any{"last_used_at": &now}).Error
}
