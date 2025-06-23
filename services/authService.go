package services

import (
	"errors"
	"fmt"
	"wellnesspath/config"
	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/models"
	"wellnesspath/repositories"
)

type AuthService struct{}

func (s *AuthService) RegisterUser(data dto.CreateUserDTO) (dto.CredentialResponseDTO, error) {
	hashedPassword, err := helpers.HashPassword(data.Password)
	if err != nil {
		return dto.CredentialResponseDTO{}, err
	}

	user := &models.User{
		Name:     data.Name,
		Username: data.Username,
		Password: hashedPassword,
		Profile:  "https://anggurproject.blob.core.windows.net/syncspend/profile/default.png",
	}

	tx := config.DB.Begin()

	if err := repositories.InsertUser(tx, user); err != nil {
		tx.Rollback()
		return dto.CredentialResponseDTO{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dto.CredentialResponseDTO{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	accessToken, err := helpers.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return dto.CredentialResponseDTO{}, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return dto.CredentialResponseDTO{}, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return dto.CredentialResponseDTO{
		ID:           user.ID,
		Name:         user.Name,
		Profile:      user.Profile,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) AuthenticateUser(data dto.LoginCredentialsDTO) (dto.CredentialResponseDTO, error) {
	tx := config.DB.Begin()

	user, err := repositories.GetUserByUsername(tx, data.Username)
	if err != nil {
		tx.Rollback()
		return dto.CredentialResponseDTO{}, errors.New("user not found")
	}

	err = helpers.CheckPasswordHash(data.Password, user.Password)
	if err != nil {
		tx.Rollback()
		return dto.CredentialResponseDTO{}, errors.New("invalid credentials")
	}

	if err := tx.Commit().Error; err != nil {
		return dto.CredentialResponseDTO{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	accessToken, err := helpers.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return dto.CredentialResponseDTO{}, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return dto.CredentialResponseDTO{}, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return dto.CredentialResponseDTO{
		ID:           user.ID,
		Name:         user.Name,
		Profile:      user.Profile,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
