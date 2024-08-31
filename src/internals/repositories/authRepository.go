package repositories

import (
	"github.com/timebetov/readerblog/internals/models"
	"github.com/timebetov/readerblog/internals/utils"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (ar *authRepository) FindUserByCredentials(username, password string) (*models.User, error) {
	var user models.User
	if err := ar.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	if err := utils.CheckPassword(user.Password, password); err != nil {
		return nil, err
	}
	return &user, nil
}
func (ar *authRepository) FindSelf(username string) (*models.User, error) {
	var user models.User
	// Find the user with the matching username
	err := ar.db.First(&user, "username = ?", username).Error
	return &user, err
}
