package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func UserToResponse(user *entity.User) *model.UserResponse {
	if user == nil {
		return nil
	}

	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  stringValue(user.Username),
		FullName:  stringValue(user.FullName),
		AvatarURL: stringValue(user.AvatarURL),
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserToAuthUser(user *entity.User) *model.AuthUser {
	if user == nil {
		return nil
	}

	return &model.AuthUser{
		ID:       user.ID,
		Email:    user.Email,
		Username: stringValue(user.Username),
		Role:     user.Role,
	}
}
