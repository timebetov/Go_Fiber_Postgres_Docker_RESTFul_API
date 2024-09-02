package services

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/timebetov/readerblog/internals/models/dtos"
	"github.com/timebetov/readerblog/internals/repositories"
	"github.com/timebetov/readerblog/internals/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	repo        repositories.AuthRepository
	userService *UserService
	redisClient *redis.Client
}

var ctx = context.Background()

func NewAuthService(repo repositories.AuthRepository, userService *UserService, redisClient *redis.Client) *AuthService {
	return &AuthService{repo, userService, redisClient}
}

func (as *AuthService) RegisterUser(userDTO *dtos.CreateUserDTO) (*dtos.ProfileDTO, string, error) {
	user, err := as.userService.CreateUser(userDTO)
	if err != nil {
		return nil, "", err
	}

	// Generating the token
	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		return nil, "", err
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

	return userDto, token, nil
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

func (as *AuthService) Logout(token string) error {
	claims, err := utils.ParseToken(token)
	if err != nil {
		return err
	}
	// Using the token's expiration time for Redis expiration
	expiration := time.Until(claims.ExpiresAt.Time)
	if expiration < 0 {
		expiration = 0
	}

	return as.redisClient.Set(ctx, token, "blacklisted", expiration).Err()
}

func (as *AuthService) IsTokenBlacklisted(token string) bool {
	_, err := as.redisClient.Get(ctx, token).Result()
	return err == nil
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
