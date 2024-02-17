package handlers

import (
	"employee-management/app/models"
	"employee-management/app/utils"
	"employee-management/config"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	expiresAccessToken = time.Now().Add(time.Minute * 1440).Unix()
)

func RegisterUserHandler(c *fiber.Ctx) error {
	user := new(models.User)
	var existingUser models.User
	// Request body check
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body tidak valid.",
		})
	}

	// Validate
	errorList := utils.ValidateRegisterUser(user)
	if errorList != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Permintaan di tolak.",
			"errors":  errorList,
		})
	}

	// Username check
	if err := config.DB.First(&existingUser, "username = ?", user.Username).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Username sudah ada.",
		})
	}

	// Email check
	if err := config.DB.First(&existingUser, "email = ?", user.Email).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email sudah ada.",
		})
	}

	// Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal hash password.",
		})
	}
	user.Password = string(hashedPassword)

	// Save to DB
	if err := config.DB.Create(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat user.",
			"error":   err.Error(),
		})
	}

	newUser := fiber.Map{
		"message": fmt.Sprintf("Selamat %s, akun anda telah dibuat.", user.FullName),
		"user": map[string]interface{}{
			"full_name":  user.FullName,
			"created_at": user.CreatedAt,
		},
	}

	// Mengembalikan data user yang berhasil dibuat
	return c.Status(fiber.StatusCreated).JSON(newUser)
}

func LoginHandlerToken(c *fiber.Ctx) error {
	// Membuat objek login dari data yang diterima
	login := new(models.User)
	if err := c.BodyParser(login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body tidak valid.",
		})
	}

	// Mengambil pengguna dengan username yang diberikan dari database
	var user models.User
	if err := config.DB.Where("username = ?", login.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Pengguna tidak ditemukan
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Username atau password salah.",
			})
		}
		// Terjadi kesalahan lain saat mengambil pengguna
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch user.",
			"error":   err.Error(),
		})
	}

	// Memverifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		// Password tidak cocok
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Username atau password salah.",
		})
	}

	// Membuat access token
	accessToken, err := createAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat aksess token.",
			"error":   err.Error(),
		})
	}

	// Membuat refresh token
	refreshToken, err := createRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat refresh token.",
			"error":   err.Error(),
		})
	}

	// Mengembalikan respons dengan token JWT
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    expiresAccessToken,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Get("Authorization")

	// Implementasi validasi dan pengambilan user berdasarkan refresh token menggunakan GORM
	user, err := getUserFromRefreshToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token tidak valid."})
	}

	// Membuat access token baru
	accessToken, err := createAccessToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat aksess token."})
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    expiresAccessToken,
	})
}

// Membuat access token
func createAccessToken(user models.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"type": "Bearer",
		"exp":  expiresAccessToken,
	})

	token, err := claims.SignedString(config.SecretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Membuat refresh token
func createRefreshToken(user models.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"type": "Refresh Token",
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token berlaku selama 1 minggu
	})

	token, err := claims.SignedString(config.SecretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func getUserFromRefreshToken(refreshToken string) (*models.User, error) {
	// Mendapatkan payload dari token
	payload, err := getPayloadFromToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Mendapatkan user ID dari payload
	userID := payload["id"].(string) // Mengubah tipe data sesuai dengan definisi model Anda

	// Mencari user berdasarkan user ID
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func getPayloadFromToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token tidak valid")
}
