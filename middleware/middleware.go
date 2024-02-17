package middleware

import (
	"employee-management/app/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func PublicMiddleware(c *fiber.Ctx) error {
	return c.Next()
}

func PrivateMiddleware(c *fiber.Ctx) error {
	// Cek apakah header Authorization ada
	tokenString := c.Get("Authorization")
	if len(tokenString) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Anda belum login!"})
	}
	// Verifikasi token JWT
	token, err := utils.VerifyToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid."})
	}

	// Periksa apakah token telah kedaluwarsa
	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token kedaluwarsa."})
	}

	// Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses klaim token."})
	}

	// Verify Bearer
	bearer := claims["type"].(string)
	if bearer != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Jenis token tidak valid."})
	}

	// Periksa apakah token memiliki waktu kedaluwarsa yang cukup
	exp := claims["exp"].(float64)
	expTime := time.Unix(int64(exp), 0)
	if time.Now().After(expTime) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token kedaluwarsa."})
	}

	return c.Next()
}

func PrivateMiddlewareAdmin(c *fiber.Ctx) error {
	// Cek apakah header Authorization ada
	tokenString := c.Get("Authorization")
	if len(tokenString) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Anda belum login!"})
	}
	// Verifikasi token JWT
	token, err := utils.VerifyToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid."})
	}

	// Periksa apakah token telah kedaluwarsa
	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token kedaluwarsa."})
	}

	// Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses klaim token."})
	}

	// Verify Bearer
	bearer := claims["type"].(string)
	if bearer != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Jenis token tidak valid."})
	}

	// Periksa apakah token memiliki waktu kedaluwarsa yang cukup
	exp := claims["exp"].(float64)
	expTime := time.Unix(int64(exp), 0)
	if time.Now().After(expTime) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token kedaluwarsa."})
	}

	// Akses admin
	userRole := claims["role"].(string)
	if userRole != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Maaf anda bukan admin."})
	}

	return c.Next()
}
