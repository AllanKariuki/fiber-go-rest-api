package routes

import (
	"github.com/AllanKariuki/fiber-go-rest-api/controllers"
	"github.com/AllanKariuki/fiber-go-rest-api/services"
	"github.com/AllanKariuki/fiber-go-rest-api/config"
	"github.com/AllanKariuki/fiber-go-rest-api/middleware"
	"github.com/AllanKariuki/fiber-go-rest-api/repositories"

	"github.com/gofiber/fiber/v2"

)

func SetupRoutes(app *fiber.App) {
	// initialize dependencies
	db := config.GetDB()

	// Repositories
	userRepo := repositories.NewUserRepository(db)

	// Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)

	// Controllers
	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// API v1 routes
	api := app.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)

	// Protected routes
	protected := api.Group("/", authMiddleware.Protected)
	protected.Get("/profile", authController.GetProfile)

	// User routes
	users := protected.Group("/users")
	users.Get("/", userController.GetAllUsers)
	users.Get("/:id", userController.GetUser)
	users.Put("/:id", userController.UpdateUser)

	// Admin only routes
	admin := users.Group("/", authMiddleware.AdminOnly)
	admin.Delete("/:id", userController.DeleteUser)
}
