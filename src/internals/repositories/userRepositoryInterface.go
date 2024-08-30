package repositories

import (
	"github.com/timebetov/readerblog/internals/models"
)

type UserRepository interface {
	FindAllUsers() ([]models.User, error)
	FindDeletedUsers() ([]models.User, error)
	FindUserById(id string) (*models.User, error)
	FindUserByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	SoftDeleteUser(user *models.User) error
	ForceDeleteUser(user *models.User) error
	RestoreUser(user *models.User) error
	SetRole(user *models.User, role string) error

	// Auth
	FindUserByCredentials(username, password string) (*models.User, error)
}
