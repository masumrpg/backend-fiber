package routes

import (
	"employee-management/app/handlers"
	"employee-management/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/api/auth")

	auth.Post("/register", middleware.PublicMiddleware, handlers.RegisterUserHandler)
	auth.Post("/token", middleware.PublicMiddleware, handlers.LoginHandlerToken) // SignIn
	auth.Post("/refresh", middleware.PrivateMiddleware, handlers.RefreshToken)
}
