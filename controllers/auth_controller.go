package controllers

import (
	"strconv"

	"github.com/AllanKariuki/fiber-go-rest-api/models"
	"github.com/AllanKariuki/fiber-go-rest-api/services"
	"github.com/AllanKariuki/fiber-go-rest-api/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthController(authService services.AuthService, userService services.UserService) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
	}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	dto := new(models.RegisterDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}
	user, token, err := ac.authService.Register(dto)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, "Register successfull", fiber.Map{
		"user": user,
		"token": token,
	})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	dto := new(models.LoginDTO)
	if err := c.BodyParser(dto); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	user, token, err := ac.authService.Login(dto)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", fiber.Map{
		"user":  user,
		"token": token,
	})
}

func (ac *AuthController) GetProfile(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Locals("user_id").(string), 10, 32)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}
	user, err := ac.userService.GetUserByID(uint(userID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "User profile retrieved successfully", user)
}