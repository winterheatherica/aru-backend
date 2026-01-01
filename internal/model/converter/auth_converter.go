package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func UserToAuthResponse(
	user *entity.User,
	accessToken string,
	refreshToken string,
) *model.AuthResponse {

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         UserToResponse(user),
	}
}
