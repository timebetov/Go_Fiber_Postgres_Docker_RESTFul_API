package controllers

import (
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
