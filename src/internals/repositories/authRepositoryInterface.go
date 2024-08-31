package repositories

import (
	"github.com/timebetov/readerblog/internals/models"
)

type AuthRepository interface {
	FindUserByCredentials(username, password string) (*models.User, error)
	FindSelf(id string) (*models.User, error)
}
