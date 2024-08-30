package services

import (
	"errors"

	"github.com/timebetov/readerblog/internals/models"
	"github.com/timebetov/readerblog/internals/repositories"
	"github.com/timebetov/readerblog/internals/utils"
)

type AuthService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) Authenticate(username, password string) (string, error) {
	user, err := s.repo.FindUserByCredentials(username, password)
	if err != nil {
		return "", errors.New("unfortunately, User not found")
	}

	// Generating JWT Token
	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GetUserProfile(username string) (*models.User, error) {
	if user, err := s.repo.FindUserByUsername(username); err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (s *AuthService) RegisterUser(user *models.User) error {
	if err := s.repo.CreateUser(user); err != nil {
		return err
	}
	return nil
}
