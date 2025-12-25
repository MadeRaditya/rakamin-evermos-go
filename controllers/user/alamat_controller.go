package user

import (
	"backend-evermos/database"
	"backend-evermos/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AlamatInput struct {
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp"`
	DetailAlamat string `json:"detail_alamat"`
}

type AlamatUpdate struct {
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp"`
	DetailAlamat string `json:"detail_alamat"`
}

func GetAllAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token tidak valid"})
	}

	var alamat []models.Alamat
	if err := database.DB.Where("id_user = ?", userID).Find(&alamat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil alamat",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": alamat,
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Alamat berhasil ditambahkan",
		"data":    alamat,
	})
}

func GetAlamatByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token tidak valid"})
	}

	var alamatID uint
	if _, err := fmt.Sscan(id, &alamatID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
		})
	}

	var alamat models.Alamat
	if err := database.DB.Where("id = ? AND id_user = ?", alamatID, userID).First(&alamat).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Alamat tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": alamat,
	})
}

func UpdateAlamat(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")
	var input AlamatUpdate

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token tidak valid"})
	}

	var alamatID uint
	if _, err := fmt.Sscan(id, &alamatID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
		})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input tidak valid",
		})
	}

	var alamat models.Alamat
	if err := database.DB.Where("id = ? AND id_user = ?", alamatID, userID).First(&alamat).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Alamat tidak ditemukan"})
	}

	alamat.NamaPenerima = input.NamaPenerima
	alamat.NoTelp = input.NoTelp
	alamat.DetailAlamat = input.DetailAlamat

	if err := database.DB.Save(&alamat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal Memperbarui alamat"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Alamat berhasil diperbarui",
		"data":    alamat,
	})
}

func DeleteAlamat(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token tidak valid"})
	}

	var alamatID uint
	if _, err := fmt.Sscan(id, &alamatID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
		})
	}

	var alamat models.Alamat
	if err := database.DB.Where("id = ? AND id_user = ?", alamatID, userID).First(&alamat).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Alamat tidak ditemukan",
		})
	}

	if err := database.DB.Delete(&alamat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus alamat",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Alamat berhasil dihapus",
		"data":    alamat,
	})
}
