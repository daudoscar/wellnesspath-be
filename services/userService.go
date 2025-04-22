package services

import (
	"errors"

	"syncspend/dto"
	"syncspend/helpers"
	"syncspend/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func (s *UserService) UpdateUser(data dto.UpdateUserDTO) (dto.CredentialResponseDTO, error) {
	user, err := repositories.GetUserByID(data.ID)
	if err != nil {
		return dto.CredentialResponseDTO{}, errors.New("user not found")
	}

	if data.Name != "" {
		user.Name = data.Name
	}

	if data.Profile != nil {
		profileImageURL, err := helpers.UploadProfileImage(data.Profile, int(user.ID))
		if err != nil {
			return dto.CredentialResponseDTO{}, err
		}

		user.Profile = profileImageURL
	}

	if data.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return dto.CredentialResponseDTO{}, errors.New("failed to hash password")
		}
		user.Password = string(hashedPassword)
	}

	if err := repositories.UpdateUser(user); err != nil {
		return dto.CredentialResponseDTO{}, errors.New("failed to update user")
	}

	userResponse := dto.CredentialResponseDTO{
		ID:       user.ID,
		Name:     user.Name,
		Profile:  user.Profile,
		Username: user.Username,
	}

	return userResponse, nil
}

func (s *UserService) GetUserByID(data dto.GetUserDTO) (dto.GetUserResponse, error) {
	user, err := repositories.GetUserByID(data.ID)
	if err != nil {
		return dto.GetUserResponse{}, errors.New("user not found")
	}

	userResponse := dto.GetUserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Profile:  user.Profile,
		Username: user.Username,
	}

	return userResponse, nil
}
