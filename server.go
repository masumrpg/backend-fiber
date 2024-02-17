package main

import (
	"employee-management/app/routes"
	"employee-management/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Setup koneksi database
	config.SetupDatabase()

	app := fiber.New(fiber.Config{})
	app.Get("/", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Runing..."}) })

	// Route app
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)

	app.Listen(":8080")
}
