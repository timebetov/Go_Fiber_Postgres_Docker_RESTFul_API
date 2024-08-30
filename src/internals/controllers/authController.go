package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/timebetov/readerblog/internals/models"
	"github.com/timebetov/readerblog/internals/services"
	"github.com/timebetov/readerblog/internals/utils"
)

type AuthController struct {
	Service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{Service: service}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username" validate:"required,username,min=8,max=32"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid input"})
	}

	// Checking the user if already exists
	if _, err := ac.Service.GetUserProfile(input.Username); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  "error",
			"message": "user already exists"})
	}

	// To lowercase the email and username
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Username = strings.ToLower(strings.TrimSpace(input.Username))

	// Validate the user data
	if err := utils.ValidateUser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error()})
	}

	// Hashing the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "could not hash password"})
	}

	// Creating a new instance of user model
	user := &models.User{
		ID:       uuid.New(),
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// Saving the user to the database
	if err := ac.Service.RegisterUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "failed to create user",
		})
	}

	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to generate token"})
	}

	// Return in success case
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User was registered successfully!",
		"data":    user,
		"token":   token,
	})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	type loginRequest struct {
		Username string `json:"username" validate:"required,username,min=8,max=32"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error()})
	}

	// Converting the username field to lowercase and trim any spaces before and after
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))

	// Validate the user data
	if err := utils.ValidateUser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"error":  err.Error(),
		})
	}

	token, err := ac.Service.Authenticate(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Authentication failed!",
			"data":    req,
			"error":   err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   token,
	})
}

func (ac *AuthController) Profile(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.Claims)
	username := claims.Username
	user, err := ac.Service.GetUserProfile(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"error":   err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}
