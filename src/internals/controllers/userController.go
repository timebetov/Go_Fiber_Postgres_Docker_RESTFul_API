package controllers

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/timebetov/readerblog/internals/models"
	"github.com/timebetov/readerblog/internals/repositories"
	"github.com/timebetov/readerblog/internals/utils"
)

type UserController struct {
	Repo repositories.UserRepository
}

func NewUserController(repo repositories.UserRepository) *UserController {
	return &UserController{
		Repo: repo,
	}
}

// Getting all users
func (userControl *UserController) GetUsers(c *fiber.Ctx) error {
	// Getting all users
	users, err := userControl.Repo.FindAllUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not retrieve users",
			"error":   err})
	}

	// If no users were found, return an error
	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No users found!",
			"data":    nil})
	}

	// In case of success, return the users if found at least 1 user
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Users found!",
		"data":    &users})
}

func (userControl *UserController) GetDeletedUsers(c *fiber.Ctx) error {
	users, err := userControl.Repo.FindDeletedUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not retrieve users!",
			"error":   err})
	}
	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No users found"})
	}

	// In case of success, return the users if found at least 1 user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Found deleted users!",
		"data":    &users,
	})
}

// Creating a brand new User
func (userControl *UserController) CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	// Parse request body into user struct
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"data":    err})
	}

	// Converting the username field to lowercase and trim any spaces before and after
	user.Username = strings.ToLower(strings.TrimSpace(user.Username))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	// Validate the user data
	if err := utils.ValidateUser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error()})
	}

	// Hashing the password before saving it to the database
	if hashedPassword, err := utils.HashPassword(user.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to hash password!"})
	} else {
		user.Password = hashedPassword
	}

	// Saving the user to the database with UUID
	user.ID = uuid.New()

	if err := userControl.Repo.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't create user",
			"data":    err})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User was created successfully!",
		"data":    user})
}

func (userControl *UserController) GetUser(c *fiber.Ctx) error {
	// Read the param userId
	id := c.Params("userId")

	// Getting the user or returning an error if not found
	user, err := userControl.Repo.FindUserById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No user found with ID"})
	}

	// In case of success return the user
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User was found successfully!",
		"data":    user})
}

func (userControl *UserController) UpdateUser(c *fiber.Ctx) error {
	// defining the struct for updating the user
	type updateUser struct {
		Email    string `json:"email" validate:"email"`
		Password string `json:"password" validate:"min=8,max=32"`
	}

	// Getting the userId from params
	id := c.Params("userId")

	// Getting the user or returning an error if not found
	user, err := userControl.Repo.FindUserById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No user found with ID"})
	}

	var updateUserData updateUser
	if err := c.BodyParser(&updateUserData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "erorr",
			"message": "Review your input"})
	}

	// Validate the user data
	if err := utils.ValidateUser(updateUserData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error()})
	}

	user.Email = updateUserData.Email

	// Hashing the password before saving it to the database
	if hashedPassword, err := utils.HashPassword(updateUserData.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to hash password!"})
	} else {
		user.Password = hashedPassword
	}

	// Saving changes
	if err = userControl.Repo.UpdateUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't update user",
			"data":    err})
	}

	// Returning the updated user
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User was updated successfully!",
		"data":    user})
}

func (userControl *UserController) SoftDeleteUser(c *fiber.Ctx) error {
	// Read the param userId
	id := c.Params("userId")

	user, err := userControl.Repo.FindUserById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No user found with ID",
			"data":    nil})
	}

	if err = userControl.Repo.SoftDeleteUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't delete user",
			"data":    err})
	}

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User: " + user.Username + " was deleted successfully"})
}

func (userControl *UserController) ForceDeleteUser(c *fiber.Ctx) error {
	// Read the param userId
	id := c.Params("userId")

	user, err := userControl.Repo.FindUserById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No user found with ID"})
	}

	if err = userControl.Repo.ForceDeleteUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't delete user",
			"data":    err})
	}

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User: " + user.Username + " was permanently deleted successfully",
	})
}

func (userControl *UserController) RestoreUser(c *fiber.Ctx) error {
	// Read the param userId
	id := c.Params("userId")

	// Getting the specified user
	user, err := userControl.Repo.FindUserById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No user found with ID"})
	}

	// Restoring the user
	err = userControl.Repo.RestoreUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't restore user",
			"data":    err})
	}

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User: " + user.Username + " was restored successfully",
	})
}

func (userControl *UserController) SetUserRole(c *fiber.Ctx) error {
	type UseRole struct {
		Role string `json:"role" validate:"required,min=5"`
	}
	// Read the param userId
	id := c.Params("userId")

	// Getting the specified user
	user, err := userControl.Repo.FindUserById(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "No user found with ID"})
	}

	var userole UseRole

	// Getting the role to set from the body
	if err := c.BodyParser(&userole); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"data":    err})
	}

	if err := utils.ValidateUser(userole); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error()})
	}

	// Check if is role matches with the allowed roles
	userole.Role = strings.ToLower(userole.Role)
	if userole.Role != os.Getenv("ADMIN_ROLE") || userole.Role != os.Getenv("READER_ROLE") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid role",
		})
	}

	// Setting the role
	err = userControl.Repo.SetRole(user, userole.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't set role",
			"data":    err})
	}

	// Return in success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Role was set successfully",
	})
}
