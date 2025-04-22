package services

import (
	"errors"
	"fmt"

	"syncspend/dto"
	"syncspend/helpers"
	"syncspend/models"
	"syncspend/repositories"
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

	if err := repositories.InsertUser(user); err != nil {
		return dto.CredentialResponseDTO{}, err
	}

	accessToken, err := helpers.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return dto.CredentialResponseDTO{}, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return dto.CredentialResponseDTO{}, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	credentialResponse := dto.CredentialResponseDTO{
		ID:           user.ID,
		Name:         user.Name,
		Profile:      user.Profile,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return credentialResponse, nil
}

func (s *AuthService) AuthenticateUser(data dto.LoginCredentialsDTO) (dto.CredentialResponseDTO, error) {
	user, err := repositories.GetUserByUsername(data.Username)
	if err != nil {
		return dto.CredentialResponseDTO{}, errors.New("user not found")
	}

	err = helpers.CheckPasswordHash(data.Password, user.Password)
	if err != nil {
		return dto.CredentialResponseDTO{}, errors.New("invalid credentials")
	}

	accessToken, err := helpers.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return dto.CredentialResponseDTO{}, fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return dto.CredentialResponseDTO{}, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	userResponse := dto.CredentialResponseDTO{
		ID:           user.ID,
		Name:         user.Name,
		Profile:      user.Profile,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return userResponse, nil
}
