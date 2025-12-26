package user

import (
	"backend-evermos/database"
	"backend-evermos/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type userInput struct {
	Nama         string `json:"nama"`
	KataSandi    string `json:"kata_sandi"`
	NoTelp       string `json:"no_telp"`
	TanggalLahir string `json:"tanggal_Lahir"`
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}

func GetMyProfil(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"Token tidak valid"},
			"data":    nil,
		})
	}

	var user models.User
	if err := database.DB.Preload("Toko").First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"User tidak ditemukan"},
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    user,
	})
}

func UpdateProfil(c *fiber.Ctx) error {
	var input userInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{"Input tidak valid"},
			"data":    nil,
		})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{"Token tidak valid"},
			"data":    nil,
		})
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{"User tidak ditemukan"},
			"data":    nil,
		})
	}

	if strings.TrimSpace(input.KataSandi) != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.KataSandi), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to PUT data",
				"errors":  []string{"Gagal meng-hash password"},
				"data":    nil,
			})
		}
		user.KataSandi = string(hashedPassword)
	}

	if strings.TrimSpace(input.TanggalLahir) != "" {
		tgl, err := time.Parse("02/01/2006", input.TanggalLahir)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to PUT data",
				"errors":  []string{"Format tanggal_lahir harus DD/MM/YYYY"},
				"data":    nil,
			})
		}
		user.TanggalLahir = tgl
	}

	if input.Nama != "" {
		user.Nama = input.Nama
	}
	if input.NoTelp != "" {
		user.NoTelp = input.NoTelp
	}
	if input.Pekerjaan != "" {
		user.Pekerjaan = input.Pekerjaan
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.IDProvinsi != "" {
		user.IDProvinsi = input.IDProvinsi
	}
	if input.IDKota != "" {
		user.IDKota = input.IDKota
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    "",
	})
}
