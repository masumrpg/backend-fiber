package config

import (
	"employee-management/app/models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	SecretKey []byte
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}
	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}

func SetupDatabase() {
	var err error
	dsn := os.Getenv("DB_URI")

	// Membuat koneksi database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	defer DB.DB()

	// Menjalankan migrasi jika diperlukan
	DB.AutoMigrate(&models.User{}, &models.UserDetail{}, &models.Address{})
}
