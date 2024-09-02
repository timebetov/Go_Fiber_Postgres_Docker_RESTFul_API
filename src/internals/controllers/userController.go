package controllers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/timebetov/readerblog/internals/models/dtos"
	"github.com/timebetov/readerblog/internals/services"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{
		Service: service,
	}
}

// Getting all users
func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	// Getting query
	deletedQuery := c.Query("deleted")

	// Getting all users
	users, err := uc.Service.GetUsers(deletedQuery)
	if err != nil {
		if err.Error() == "invalid deleted query" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid boolean value for 'deleted' query parameter",
			})
		}
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

// Creating a brand new User
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	var userDTO dtos.CreateUserDTO

	// Parse request body into userDTO struct
	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Review your input",
			"data":    err})
	}

	// Passing to the service layer to create a new user
	createdUser, err := uc.Service.CreateUser(&userDTO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't create user",
			"data":    err.Error(),
		})
	}

	// Return in success case
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User was registered successfully!",
		"data":    createdUser,
	})
}

func (uc *UserController) GetUser(c *fiber.Ctx) error {
	// Read the param userId
	id := c.Params("userId")

	// Getting the user or returning an error if not found
	user, err := uc.Service.GetUserById(id)
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

func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	var userDTO dtos.UpdateUserDTO

	// Getting the userId from params
	id := c.Params("userId")

	// Parsing the request body into userDTO
	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "erorr",
			"message": "Review your input"})
	}

	user, err := uc.Service.UpdateUser(id, &userDTO)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "No user found with ID",
			})
		} else if err.Error() == "invalid role" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid role",
				"allowed": []string{os.Getenv("ADMIN_ROLE"), os.Getenv("WRITER_ROLE")},
			})
		} else if err.Error() == "passwords do not match" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Passwords do not match",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	// Returning the updated user
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User was updated successfully!",
		"data":    user})
}

func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	// Read the param userId
	id := c.Params("userId")
	// Read the query 'force'
	forceQuery := c.Query("force")

	user, err := uc.Service.DeleteUser(forceQuery, id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "No user found with ID"})
		} else if err.Error() == "invalid force query" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid boolean value for 'force' query parameter",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't delete user",
			"error":   err.Error(),
		})
	}

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User: " + user.Username + " was deleted successfully"})
}

func (uc *UserController) RestoreUser(c *fiber.Ctx) error {
	// Read the param userId
	id := c.Params("userId")

	// Getting the specified user
	user, err := uc.Service.RestoreUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "No user found with ID",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't restore user",
			"data":    err.Error(),
		})
	}

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User: " + user.Username + " was restored successfully",
	})
}
