package services

import (
	"errors"
	"wellnesspath/config"
	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func (s *UserService) UpdateUser(data dto.UpdateUserDTO) (dto.CredentialResponseDTO, error) {
	tx := config.DB.Begin()

	user, err := repositories.GetUserByID(tx, data.ID)
	if err != nil {
		tx.Rollback()
		return dto.CredentialResponseDTO{}, errors.New("user not found")
	}

	if data.Name != "" {
		user.Name = data.Name
	}

	if data.Profile != nil {
		profileImageURL, err := helpers.UploadProfileImage(data.Profile, int(user.ID))
		if err != nil {
			tx.Rollback()
			return dto.CredentialResponseDTO{}, err
		}
		user.Profile = profileImageURL
	}

	if data.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			return dto.CredentialResponseDTO{}, errors.New("failed to hash password")
		}
		user.Password = string(hashedPassword)
	}

	if err := repositories.UpdateUser(tx, user); err != nil {
		tx.Rollback()
		return dto.CredentialResponseDTO{}, errors.New("failed to update user")
	}

	if err := tx.Commit().Error; err != nil {
		return dto.CredentialResponseDTO{}, errors.New("failed to commit transaction")
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
	tx := config.DB.Begin()

	user, err := repositories.GetUserByID(tx, data.ID)
	if err != nil {
		tx.Rollback()
		return dto.GetUserResponse{}, errors.New("user not found")
	}

	if err := tx.Commit().Error; err != nil {
		return dto.GetUserResponse{}, errors.New("failed to commit transaction")
	}

	userResponse := dto.GetUserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Profile:  user.Profile,
		Username: user.Username,
	}

	return userResponse, nil
}

func (s *UserService) DeleteUser(data dto.DeleteUserDTO) error {
	tx := config.DB.Begin()

	if err := repositories.DeleteUserByID(tx, data.ID); err != nil {
		tx.Rollback()
		return errors.New("failed to delete user")
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction")
	}

	return nil
}
