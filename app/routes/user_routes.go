package routes

import (
	"employee-management/app/handlers"
	"employee-management/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	user := app.Group("/api/user")
	users := app.Group("/api/users")

	user.Post("/", middleware.PrivateMiddleware, handlers.CreateUserDetailHandler)
	user.Get("/me", middleware.PrivateMiddleware, handlers.GetMeHandler)
	user.Get("/:id", middleware.PrivateMiddleware, handlers.GetUserByIDHandler)
	user.Get("/:id/detailed", middleware.PrivateMiddleware, handlers.GetUserByIDDetailedHandler)
	user.Put("/:id/detailed", middleware.PrivateMiddleware, handlers.UpdateUserWithDetailHandler)
	user.Delete("/:id", middleware.PrivateMiddleware, handlers.DeleteUserHandler)

	users.Get("/", middleware.PrivateMiddlewareAdmin, handlers.GetAllUsersHandler)
}
