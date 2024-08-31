package dtos

import "time"

type CreateUserDTO struct {
	Username             string `json:"username" validate:"required,min=8,max=32,username"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,max=32"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type UpdateUserDTO struct {
	Email                *string `json:"email" validate:"omitempty,email"`
	Password             *string `json:"password" validate:"omitempty,min=8,max=32"`
	PasswordConfirmation *string `json:"password_confirmation" validate:"omitempty,eqfield=Password"`
	Role                 *string `json:"role" validate:"omitempty,min=5"`
}

type LoginUserDTO struct {
	Username             string `json:"username" validate:"required,username,min=8,max=32"`
	Password             string `json:"password" validate:"required,min=8,max=32"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type ProfileDTO struct {
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Subscribers uint      `json:"subscribers"`
	Followed    uint      `json:"followed"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
}
