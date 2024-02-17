package utils

import (
	"employee-management/config"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetIDFromAccessToken(c *fiber.Ctx) (string, error) {
	// Dapatkan nilai token JWT dari header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("token not provided")
	}

	// Pisahkan header Authorization untuk mendapatkan token JWT
	auth := strings.Split(authHeader, " ")
	if len(auth) == 0 {
		return "", errors.New("invalid authorization header format")
	}
	tokenString := auth[0]

	// Parse token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Anda perlu menyesuaikan dengan metode verifikasi token Anda, misalnya, menggunakan secret key
		return config.SecretKey, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid access_token")
	}

	// Ambil klaim ID dari token JWT
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("failed to parse token claims")
	}

	userID := claims["id"].(string)

	return userID, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Parse token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Anda perlu menyesuaikan dengan metode verifikasi token Anda, misalnya, menggunakan secret key
		return config.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
