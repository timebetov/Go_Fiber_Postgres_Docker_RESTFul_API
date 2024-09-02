package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/timebetov/readerblog/internals/models"
	"github.com/timebetov/readerblog/internals/models/dtos"
	"github.com/timebetov/readerblog/internals/repositories"
	"github.com/timebetov/readerblog/internals/utils"
	"gorm.io/gorm"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo}
}

func (us *UserService) GetUserById(id string) (*models.User, error) {
	user, err := us.repo.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetUsers(deletedQuery string) ([]models.User, error) {
	var deleted bool

	if deletedQuery != "" {
		var err error
		// Converting the string to a boolean if provided
		deleted, err = strconv.ParseBool(deletedQuery)
		if err != nil {
			return nil, errors.New("invalid deleted query")
		}
	} else {
		// If query parameter is not provided, setting the default value
		deleted = false
	}

	users, err := us.repo.FindUsers(deleted)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *UserService) CreateUser(userDTO *dtos.CreateUserDTO) (*models.User, error) {
	// Converting the username field to lowercase and trim any spaces before and after
	userDTO.Username = utils.TrimAndLower(userDTO.Username)
	userDTO.Email = utils.TrimAndLower(userDTO.Email)

	// Validating user data
	if err := utils.ValidateUser(userDTO); err != nil {
		return nil, err
	}

	// Hashing the password
	hashedPassword, err := utils.HashPassword(userDTO.Password)
	if err != nil {
		return nil, err
	}

	// Creating a new user instance
	user := &models.User{
		ID:       uuid.New(),
		Username: userDTO.Username,
		Email:    userDTO.Email,
		Password: hashedPassword,
	}

	if err := us.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
func (us *UserService) UpdateUser(id string, userDTO *dtos.UpdateUserDTO) (*models.User, error) {
	// Fetching the user from the database
	user, err := us.repo.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Validating the DTO
	if err := utils.ValidateUser(userDTO); err != nil {
		return nil, err
	}

	// Updating the field if they are not nil
	if userDTO.Email != nil {
		user.Email = utils.TrimAndLower(*userDTO.Email)
	}
	if userDTO.Role != nil {
		newRole := utils.TrimAndLower(*userDTO.Role)
		if newRole != os.Getenv("ADMIN_ROLE") && newRole != os.Getenv("WRITER_ROLE") {
			return nil, errors.New("invalid role")
		}
		user.Role = newRole
	}
	if userDTO.Password != nil {
		if userDTO.PasswordConfirmation == nil || *userDTO.PasswordConfirmation != *userDTO.Password {
			return nil, errors.New("passwords do not match")
		}
		hashedPassword, err := utils.HashPassword(*userDTO.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	// Saving the updated user to the database
	if err := us.repo.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) DeleteUser(forceQuery string, id string) (*models.User, error) {
	user, err := us.GetUserById(id)
	if err != nil {
		return nil, err
	}
	var force bool

	if forceQuery != "" {
		var err error
		// Converting the string to a boolean if provided
		force, err = strconv.ParseBool(forceQuery)
		if err != nil {
			return nil, errors.New("invalid force query")
		}
	} else {
		// If query parameter is not provided, setting the default value
		force = false
	}

	if err := us.repo.DeleteUser(force, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) RestoreUser(id string) (*models.User, error) {
	user, err := us.repo.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Restoring the user
	if err := us.repo.RestoreUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
