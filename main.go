package main

import (
	"employee-management/app/routes"
	"employee-management/config"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.SetupDatabase()

	app := fiber.New()
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })
	app.Get("/", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Runing..."}) })

	// Route app
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}
