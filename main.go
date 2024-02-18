package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

// func main() {
// 	// ==========================================================================================
// 	config.SetupDatabase()

// 	app := fiber.New()
// 	app.Get("/", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Runing..."}) })

// 	// Route app
// 	routes.SetupAuthRoutes(app)
// 	routes.SetupUserRoutes(app)

// 	// Ambil host dari variabel lingkungan
// 	host := os.Getenv("RAILWAY_APP_HOST")
// 	if host == "" {
// 		host = "localhost" // Jika host tidak tersedia, gunakan localhost
// 	}

// 	// Tetapkan port default HTTP atau HTTPS tergantung pada mode kerja aplikasi
// 	var port string
// 	if os.Getenv("RAILWAY_APP_PORT") == "443" {
// 		port = "" // Port default HTTPS tidak perlu disertakan dalam alamat
// 	} else {
// 		port = ":8080" // Port default HTTP
// 	}

// 	// Mulai aplikasi di alamat yang ditentukan
// 	app.Listen(host + port)
// }

func main() {
	app := fiber.New()

	// Mengambil nilai PORT dari environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Port default jika PORT tidak diset
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Mendengarkan port yang disediakan oleh Railway
	app.Listen(":" + port)
}
