package services

import (
	"errors"

	"github.com/timebetov/readerblog/internals/models/dtos"
	"github.com/timebetov/readerblog/internals/repositories"
	"github.com/timebetov/readerblog/internals/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	repo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) *AuthService {
	return &AuthService{repo}
}

func (as *AuthService) Authenticate(userDTO *dtos.LoginUserDTO) (string, error) {
	// Converting the username field to lowercase and trim any spaces before and after
	userDTO.Username = utils.TrimAndLower(userDTO.Username)

	// Validate the user data
	if err := utils.ValidateUser(userDTO); err != nil {
		return "", err
	}

	user, err := as.repo.FindUserByCredentials(userDTO.Username, userDTO.Password)
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

func (as *AuthService) GetUserProfile(username string) (*dtos.ProfileDTO, error) {
	user, err := as.repo.FindSelf(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	userDto := &dtos.ProfileDTO{
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		Subscribers: user.Subscribers,
		Followed:    user.Followed,
		Image:       user.Image,
		CreatedAt:   user.CreatedAt,
	}

	return userDto, nil
}
