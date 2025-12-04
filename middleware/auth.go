package middleware

import (
	"strings"
	"github.com/AllanKariuki/fiber-go-rest-api/config"
	"github.com/AllanKariuki/fiber-go-rest-api/models"
	"github.com/AllanKariuki/fiber-go-rest-api/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	authService services.AuthService
}

func NewAuthMiddleware(authService services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (am *AuthMiddleware) Protected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": "Missing authorization header",
		})
	}

	// Extract token from "Bearer <token>"
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	token, err := am.authService.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": "Invalid or expired token",
		})
	}

	// Extract user claims
	claims := token.Claims.(jwt.MapClaims)

	// Get user from database
	db := config.GetDB()
	var user models.User
	if err := db.First(&user, claims["user_id"]).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": "User not found",
		})
	}

	// Store user in context
	c.Locals("user", &user)
	c.Locals("user_id", user.ID)
	c.Locals("role", user.Role)

	return c.Next()
}

func (am *AuthMiddleware) AdminOnly(c *fiber.Ctx) error {
	role := c.Locals("role").(string) 
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error": "Admin access required",
		})
	}
	return c.Next()
}