package handlers

import (
	"employee-management/app/models"
	"employee-management/app/utils"
	"employee-management/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateUserDetailHandler(c *fiber.Ctx) error {
	userID, err := utils.GetIDFromAccessToken(c)
	// Mengecek jika terjadi kesalahan
	if err != nil {
		// Mengakses pesan kesalahan menggunakan metode Error()
		errorMessage := err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errorMessage,
		})
	}

	// Parse request body into UserDetail struct
	var userDetails = new(models.UserDetail)
	if err := c.BodyParser(&userDetails); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body.",
			"error":   err.Error(),
		})
	}

	// Check existing
	var userDetailsExist *models.UserDetail
	if err := config.DB.First(&userDetailsExist, "user_id = ?", userID).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User detail sudah ada, silahkan gunakan metode update.",
		})
	}

	// Insert user details into database
	userDetails.UserID = userID
	if err := config.DB.Create(&userDetails).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat user details.",
			"error":   err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Berhasil menambahkan detail user.",
		"user_details": userDetails,
	})
}

func GetMeHandler(c *fiber.Ctx) error {
	var existingUser models.User
	userID, err := utils.GetIDFromAccessToken(c)
	// Mengecek jika terjadi kesalahan
	if err != nil {
		// Mengakses pesan kesalahan menggunakan metode Error()
		errorMessage := err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errorMessage,
		})
	}

	if err := config.DB.Preload("UserDetail").Preload("UserDetail.Address").Where("id = ?", userID).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Pengguna tidak ditemukan
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Login terlebih dahulu.",
			})
		}
		// Terjadi kesalahan lain saat mengambil pengguna
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch user.",
			"error":   err.Error(),
		})
	}

	responseUser := models.UserResponseWithDetailed{
		ID:            existingUser.ID,
		Username:      existingUser.Username,
		Email:         existingUser.Email,
		FullName:      existingUser.FullName,
		IsActive:      existingUser.IsActive,
		Role:          existingUser.Role,
		IsVerified:    existingUser.IsVerified,
		VerifiedAt:    existingUser.VerifiedAt,
		LastEntryDate: existingUser.LastEntryDate,
		CreatedAt:     existingUser.CreatedAt,
		UpdatedAt:     existingUser.UpdatedAt,
		UserDetailResponse: models.UserDetailResponse{
			IDCard: existingUser.UserDetail.IDCard,
			AddressResponse: models.AddressResponse{
				PostalCode:  existingUser.UserDetail.Address.PostalCode,
				Village:     existingUser.UserDetail.Address.Village,
				Subdistrict: existingUser.UserDetail.Address.Subdistrict,
				City:        existingUser.UserDetail.Address.City,
				Province:    existingUser.UserDetail.Address.Province,
				Country:     existingUser.UserDetail.Address.Country,
			},
			Phone:             existingUser.UserDetail.Phone,
			DOB:               existingUser.UserDetail.DOB,
			Gender:            existingUser.UserDetail.Gender,
			MaritalStatus:     existingUser.UserDetail.MaritalStatus,
			Religion:          existingUser.UserDetail.Religion,
			TertiaryEducation: existingUser.UserDetail.TertiaryEducation,
			Job:               existingUser.UserDetail.Job,
			Salary:            existingUser.UserDetail.Salary,
		},
	}

	return c.Status(fiber.StatusOK).JSON(responseUser)
}

func GetUserByIDHandler(c *fiber.Ctx) error {
	// Retrieve user ID from request params
	id := c.Params("id")
	var existingUser models.User
	if err := config.DB.Where("id = ?", id).First(&existingUser).Error; err != nil {
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

	responseUser := models.UserResponse{
		ID:            existingUser.ID,
		Username:      existingUser.Username,
		Email:         existingUser.Email,
		FullName:      existingUser.FullName,
		IsActive:      existingUser.IsActive,
		Role:          existingUser.Role,
		IsVerified:    existingUser.IsVerified,
		VerifiedAt:    existingUser.VerifiedAt,
		LastEntryDate: existingUser.LastEntryDate,
		CreatedAt:     existingUser.CreatedAt,
		UpdatedAt:     existingUser.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(responseUser)
}

func GetUserByIDDetailedHandler(c *fiber.Ctx) error {
	// Retrieve user ID from request params
	id := c.Params("id")
	var existingUser models.User
	if err := config.DB.Preload("UserDetail").Preload("UserDetail.Address").Where("id = ?", id).First(&existingUser).Error; err != nil {
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

	responseUser := models.UserResponseWithDetailed{
		ID:            existingUser.ID,
		Username:      existingUser.Username,
		Email:         existingUser.Email,
		FullName:      existingUser.FullName,
		IsActive:      existingUser.IsActive,
		Role:          existingUser.Role,
		IsVerified:    existingUser.IsVerified,
		VerifiedAt:    existingUser.VerifiedAt,
		LastEntryDate: existingUser.LastEntryDate,
		CreatedAt:     existingUser.CreatedAt,
		UpdatedAt:     existingUser.UpdatedAt,
		UserDetailResponse: models.UserDetailResponse{
			IDCard: existingUser.UserDetail.IDCard,
			AddressResponse: models.AddressResponse{
				PostalCode:  existingUser.UserDetail.Address.PostalCode,
				Village:     existingUser.UserDetail.Address.Village,
				Subdistrict: existingUser.UserDetail.Address.Subdistrict,
				City:        existingUser.UserDetail.Address.City,
				Province:    existingUser.UserDetail.Address.Province,
				Country:     existingUser.UserDetail.Address.Country,
			},
			Phone:             existingUser.UserDetail.Phone,
			DOB:               existingUser.UserDetail.DOB,
			Gender:            existingUser.UserDetail.Gender,
			MaritalStatus:     existingUser.UserDetail.MaritalStatus,
			Religion:          existingUser.UserDetail.Religion,
			TertiaryEducation: existingUser.UserDetail.TertiaryEducation,
			Job:               existingUser.UserDetail.Job,
			Salary:            existingUser.UserDetail.Salary,
		},
	}

	return c.Status(fiber.StatusOK).JSON(responseUser)
}

func UpdateUserWithDetailHandler(c *fiber.Ctx) error {
	// Dapatkan ID pengguna dari parameter URL
	id := c.Params("id")

	// Parse data baru dari body permintaan
	var userData models.User
	if err := c.BodyParser(&userData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse request body"})
	}

	// Cari pengguna berdasarkan ID
	var existingUser models.User
	if err := config.DB.Preload("UserDetail").Preload("UserDetail.Address").First(&existingUser, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	// Perbarui data pengguna dengan data baru
	userData.Role = "user"
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&existingUser).Updates(userData).Error; err != nil {
			return err
		}

		// Perbarui detail pengguna
		if err := tx.Model(&existingUser.UserDetail).Updates(userData.UserDetail).Error; err != nil {
			return err
		}

		// Perbarui alamat pengguna
		if err := tx.Model(&existingUser.UserDetail.Address).Updates(userData.UserDetail.Address).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update user"})
	}

	// Berhasil memperbarui data pengguna
	return c.JSON(fiber.Map{"message": "user updated successfully"})
}

// DeleteUserHandler menghapus pengguna berdasarkan ID yang diberikan
func DeleteUserHandler(c *fiber.Ctx) error {
	// Dapatkan ID pengguna dari parameter URL
	userID := c.Params("id")

	var userDetails models.UserDetail

	// Cek addres ID
	if err := config.DB.Where("user_id = ?", userID).First(&userDetails).Error; err != nil {
		// Jika user_detail tidak ditemukan, kirim respons tidak ditemukan
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User detail not found",
			})
		}

		// Jika terjadi kesalahan lain saat mengambil user_detail, kirim respons kesalahan server
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get user detail",
			"error":   err.Error(),
		})
	}

	// Membuka transaksi basis data
	tx := config.DB.Begin()

	// Hapus alamat terlebih dahulu yang terkait dengan pengguna
	if err := tx.Where("user_detail_id = ?", userDetails.ID).Delete(&models.Address{}).Error; err != nil {
		// Jika terjadi kesalahan saat menghapus alamat
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete address",
			"error":   err.Error(),
		})
	}

	// Hapus detail pengguna yang terkait dengan pengguna
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserDetail{}).Error; err != nil {
		// Jika terjadi kesalahan saat menghapus detail pengguna
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user detail",
			"error":   err.Error(),
		})
	}

	// Hapus pengguna dari basis data
	if err := tx.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
		// Jika terjadi kesalahan saat menghapus pengguna
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	// Commit transaksi jika semua operasi berhasil
	tx.Commit()

	// Mengembalikan respons berhasil jika pengguna berhasil dihapus
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Suksess menghapus user",
	})
}

func GetAllUsersHandler(c *fiber.Ctx) error {
	// Inisialisasi slice untuk menampung semua pengguna beserta detail dan alamat
	var users []models.User

	// Ambil semua pengguna beserta detail dan alamat dari basis data
	if err := config.DB.Preload("UserDetail").Preload("UserDetail.Address").Find(&users).Error; err != nil {
		// Jika terjadi kesalahan saat mengambil pengguna, kirim respons kesalahan server
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get users",
			"error":   err.Error(),
		})
	}

	var usersResponse []models.UserResponseWithDetailed
	for _, user := range users {
		// Remapping data pengguna ke UserResponse tanpa menyertakan password
		userResponse := models.UserResponseWithDetailed{
			ID:            user.ID,
			Username:      user.Username,
			Email:         user.Email,
			FullName:      user.FullName,
			IsActive:      user.IsActive,
			Role:          user.Role,
			IsVerified:    user.IsVerified,
			VerifiedAt:    user.VerifiedAt,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
			LastEntryDate: user.LastEntryDate,
			UserDetailResponse: models.UserDetailResponse{
				IDCard: user.UserDetail.IDCard,
				AddressResponse: models.AddressResponse{
					PostalCode:  user.UserDetail.Address.PostalCode,
					Village:     user.UserDetail.Address.Village,
					Subdistrict: user.UserDetail.Address.Subdistrict,
					City:        user.UserDetail.Address.City,
					Province:    user.UserDetail.Address.Province,
					Country:     user.UserDetail.Address.Country,
				},
				Phone:             user.UserDetail.Phone,
				DOB:               user.UserDetail.DOB,
				Gender:            user.UserDetail.Gender,
				MaritalStatus:     user.UserDetail.MaritalStatus,
				Religion:          user.UserDetail.Religion,
				TertiaryEducation: user.UserDetail.TertiaryEducation,
				Job:               user.UserDetail.Job,
				Salary:            user.UserDetail.Salary,
			},
		}
		usersResponse = append(usersResponse, userResponse)
	}

	// Jika pengguna berhasil ditemukan, kirimkan daftar pengguna beserta detail dan alamat sebagai respons
	return c.Status(fiber.StatusOK).JSON(usersResponse)
}
