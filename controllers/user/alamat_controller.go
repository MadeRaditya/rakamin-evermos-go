package user

import (
	"backend-evermos/database"
	"backend-evermos/models"

	"github.com/gofiber/fiber/v2"
)

type AlamatInput struct {
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp"`
	DetailAlamat string `json:"detail_alamat"`
}

func GetAllAlamat(c *fiber.Ctx) error {
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
	database.DB.Preload("Alamat").First(&result, user.ID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result.Alamat,
	})
}

func CreateAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	var input AlamatInput

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token tidak valid"})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	uid := uint(userID.(float64))
	alamat := models.Alamat{
		IDUser:       uid,
		JudulAlamat:  input.JudulAlamat,
		NamaPenerima: input.NamaPenerima,
		NoTelp:       input.NoTelp,
		DetailAlamat: input.DetailAlamat,
	}

	if err := database.DB.Create(&alamat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambahkan Alamat.",
		})
	}

	var result models.User
	database.DB.Preload("Alamat").First(&result, userID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Alamat berhasil ditambahkan",
		"data":    result.Alamat,
	})
}
