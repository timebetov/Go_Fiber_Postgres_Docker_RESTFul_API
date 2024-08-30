package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uuid.UUID      `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Username    string         `gorm:"unique;not null" json:"username" validate:"required,min=8,max=32,username"`
	Email       string         `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password    string         `gorm:"not null" json:"password" validate:"required,min=8,max=32"`
	Role        string         `gorm:"not null;default:'writer'"`
	Subscribers uint           `gorm:"default:0"`
	Followed    uint           `gorm:"default:0"`
	Image       string         `gorm:"type:text"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
