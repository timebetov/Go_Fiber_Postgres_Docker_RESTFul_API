package repositories

import (
	"github.com/timebetov/readerblog/internals/models"
)

type UserRepository interface {
	FindUsers(includeDeleted bool) ([]models.User, error)
	FindUserById(id string) (*models.User, error)
	FindUserByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(force bool, user *models.User) error
	RestoreUser(user *models.User) error
}
