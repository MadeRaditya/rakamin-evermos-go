package controllers

import (
	"backend-evermos/database"
	"backend-evermos/models"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAllToko(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	// nama     := c.Query("nama")

	offset := (page - 1) * limit

	var toko []models.Toko
	query := database.DB.Model(&models.Toko{})

	// if nama != "" {
	// 	query = query.Where("nama_toko LIKE ?", "%"+nama+"%")
	// }

	if err := query.
		Limit(limit).
		Offset(offset).
		Find(&toko).Error; err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data toko",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  toko,
		"page":  page,
		"limit": limit,
	})
}

func GetTokoByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token tidak valid"})
	}

	var tokoID uint
	if _, err := fmt.Sscan(id, &tokoID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
		})
	}

	var toko models.Toko
	if err := database.DB.Where("id = ?", tokoID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Toko tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": toko,
	})
}

func GetMyToko(c *fiber.Ctx) error {
	userID := c.Locals("user_id")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token tidak valid"})
	}

	var toko models.Toko
	if err := database.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Toko tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": toko,
	})
}

func UpdateTokoByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid",
		})
	}

	idParam := c.Params("id_toko")
	idToko, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID toko tidak valid",
		})
	}

	var toko models.Toko
	if err := database.DB.
		Where("id = ? AND id_user = ?", idToko, userID).
		First(&toko).Error; err != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Toko tidak ditemukan",
		})
	}

	namaToko := c.FormValue("nama_toko")
	if namaToko != "" {
		toko.NamaToko = namaToko
	}

	file, err := c.FormFile("photo")
	if err == nil {
		uploadDir := "./uploads/toko"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			_ = os.MkdirAll(uploadDir, os.ModePerm)
		}

		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf(
			"toko_%d_%d%s",
			toko.ID,
			time.Now().Unix(),
			ext,
		)

		filePath := filepath.Join(uploadDir, filename)

		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal menyimpan foto",
			})
		}

		toko.UrlFoto = filePath
	}

	if err := database.DB.Save(&toko).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal update toko",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Toko berhasil diupdate",
		"data":    toko,
	})
}
