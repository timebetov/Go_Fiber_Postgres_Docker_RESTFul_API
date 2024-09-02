package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/timebetov/readerblog/internals/models/dtos"
	"github.com/timebetov/readerblog/internals/services"
	"github.com/timebetov/readerblog/internals/utils"
)

type AuthController struct {
	Service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{Service: service}
}

func (ac *AuthController) RegisterUser(c *fiber.Ctx) error {
	var userDTO dtos.CreateUserDTO

	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error()})
	}

	user, token, err := ac.Service.RegisterUser(&userDTO)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Registration failed!",
			"data":    userDTO,
			"error":   err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User registered successfully!",
		"data":    user,
		"token":   token,
	})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var userDTO dtos.LoginUserDTO

	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error()})
	}

	token, err := ac.Service.Authenticate(&userDTO)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Authentication failed!",
			"data":    userDTO,
			"error":   err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"token":  token,
	})
}

func (ac *AuthController) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing or invalid token",
		})
	}

	tokenString := strings.Split(authHeader, " ")[1]
	err := ac.Service.Logout(tokenString)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to blacklist token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully logged out",
	})
}

func (ac *AuthController) Profile(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.Claims)
	username := claims.Username
	user, err := ac.Service.GetUserProfile(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "You are not authorized to access this resource"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}
