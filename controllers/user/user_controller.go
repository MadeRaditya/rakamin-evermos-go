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
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token tidak valid"})
	}

	var user models.User

	if err := database.DB.Find(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User tidak ditemukan",
		})
	}

	var result models.User
	database.DB.Preload("Toko").First(&result, user.ID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}

func UpdateProfil(c *fiber.Ctx) error {
	var input userInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid",
		})
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	if strings.TrimSpace(input.KataSandi) != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.KataSandi), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal meng-hash password"})
		}
		user.KataSandi = string(hashedPassword)
	}

	if strings.TrimSpace(input.TanggalLahir) != "" {
		tgl, err := time.Parse("02/01/2006", input.TanggalLahir)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Format tanggal_lahir harus DD/MM/YYYY"})
		}
		user.TanggalLahir = tgl
	}

	user.Nama = input.Nama
	user.NoTelp = input.NoTelp
	user.Pekerjaan = input.Pekerjaan
	user.Email = input.Email
	user.IDProvinsi = input.IDProvinsi
	user.IDKota = input.IDKota

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui profil"})
	}

	var result models.User
	database.DB.Preload("Toko").First(&result, user.ID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profil berhasil diperbarui",
		"data":    result,
	})
}
