package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	User struct {
		ID            string     `json:"id" gorm:"type:varchar(72);primaryKey"`
		Username      string     `json:"username" gorm:"unique;type:varchar(20);not null" validate:"required,min=5,max=20"`
		Password      string     `json:"password" gorm:"not null" validate:"required"`
		Email         string     `json:"email" gorm:"unique;not null" validate:"required,email"`
		FullName      string     `json:"full_name" gorm:"type:varchar(100);not null" validate:"required,min=3,max=100"`
		IsActive      bool       `json:"is_active" gorm:"default:true"`
		Role          string     `json:"role" gorm:"default:user"`
		IsVerified    bool       `json:"is_verified" gorm:"default:false"`
		VerifiedAt    *time.Time `json:"verified_at"`
		LastEntryDate *time.Time `json:"last_entry_date"`
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     time.Time  `json:"updated_at"`
		UserDetail    UserDetail `json:"user_detail"`
	}

	UserDetail struct {
		ID                uint       `json:"user_detail_id" gorm:"primaryKey"`
		Address           Address    `json:"address"`
		Phone             string     `json:"phone"`
		DOB               *time.Time `json:"dob"`
		Gender            string     `json:"gender"`
		MaritalStatus     string     `json:"marital_status"`
		IDCard            string     `json:"id_card"`
		Religion          string     `json:"religion"`
		TertiaryEducation string     `json:"tertiary_education"`
		Job               string     `json:"job"`
		Salary            int32      `json:"salary"`
		UserID            string     `json:"user_id" gorm:"foreignKey:UserID"`
	}

	Address struct {
		ID           uint   `json:"address_id" gorm:"primaryKey"`
		PostalCode   string `json:"postal_code"`
		Village      string `json:"village"`
		Subdistrict  string `json:"subdistrict"`
		City         string `json:"city"`
		Province     string `json:"province"`
		Country      string `json:"country"`
		UserDetailID string `json:"user_detail_id" gorm:"foreignKey:UserDetailID"`
	}
)

type (
	UserResponse struct {
		ID            string     `json:"id"`
		Username      string     `json:"username"`
		Email         string     `json:"email"`
		FullName      string     `json:"full_name"`
		IsActive      bool       `json:"is_active"`
		Role          string     `json:"role"`
		IsVerified    bool       `json:"is_verified"`
		VerifiedAt    *time.Time `json:"verified_at"`
		LastEntryDate *time.Time `json:"last_entry_date"`
		CreatedAt     time.Time  `json:"created_at"`
		UpdatedAt     time.Time  `json:"updated_at"`
	}

	UserResponseWithDetailed struct {
		ID                 string             `json:"id"`
		Username           string             `json:"username"`
		Email              string             `json:"email"`
		FullName           string             `json:"full_name"`
		IsActive           bool               `json:"is_active"`
		Role               string             `json:"role"`
		IsVerified         bool               `json:"is_verified"`
		VerifiedAt         *time.Time         `json:"verified_at"`
		LastEntryDate      *time.Time         `json:"last_entry_date"`
		CreatedAt          time.Time          `json:"created_at"`
		UpdatedAt          time.Time          `json:"updated_at"`
		UserDetailResponse UserDetailResponse `json:"user_detail"`
	}

	UserDetailResponse struct {
		IDCard            string          `json:"id_card"`
		AddressResponse   AddressResponse `json:"address"`
		Phone             string          `json:"phone"`
		DOB               *time.Time      `json:"dob"`
		Gender            string          `json:"gender"`
		MaritalStatus     string          `json:"marital_status"`
		Religion          string          `json:"religion"`
		TertiaryEducation string          `json:"tertiary_education"`
		Job               string          `json:"job"`
		Salary            int32           `json:"salary"`
	}

	AddressResponse struct {
		PostalCode  string `json:"postal_code"`
		Village     string `json:"village"`
		Subdistrict string `json:"subdistrict"`
		City        string `json:"city"`
		Province    string `json:"province"`
		Country     string `json:"country"`
	}
)

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return nil
}
