package controllers

import (
	"github.com/MadeRaditya/rakamin-evermos-go/database"
	"github.com/MadeRaditya/rakamin-evermos-go/models"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAllToko(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	nama := c.Query("nama")

	offset := (page - 1) * limit

	var toko []models.Toko
	query := database.DB.Model(&models.Toko{})

	if nama != "" {
		query = query.Where("nama_toko LIKE ?", "%"+nama+"%")
	}

	if err := query.Limit(limit).Offset(offset).Find(&toko).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data": fiber.Map{
			"page":  page,
			"limit": limit,
			"data":  toko,
		},
	})
}

func GetTokoByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var tokoID uint
	if _, err := fmt.Sscan(id, &tokoID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"ID tidak valid"},
			"data":    nil,
		})
	}

	var toko models.Toko
	if err := database.DB.Where("id = ?", tokoID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"Toko tidak ditemukan"},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    toko,
	})
}

func GetMyToko(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"Token tidak valid"},
			"data":    nil,
		})
	}

	var toko models.Toko
	if err := database.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"Toko tidak ditemukan"},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    toko,
	})
}

func UpdateTokoByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to UPDATE data",
			"errors":  []string{"Token tidak valid"},
			"data":    nil,
		})
	}

	idParam := c.Params("id_toko")
	idToko, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to UPDATE data",
			"errors":  []string{"ID toko tidak valid"},
			"data":    nil,
		})
	}

	var toko models.Toko
	if err := database.DB.
		Where("id = ? AND id_user = ?", idToko, userID).
		First(&toko).Error; err != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"Toko tidak ditemukan"},
			"data":    nil,
		})
	}

	namaToko := c.FormValue("nama_toko")
	if namaToko != "" {
		toko.NamaToko = namaToko
	}

	file, err := c.FormFile("photo")
	if err == nil {
		uploadDir := "./public/uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			_ = os.MkdirAll(uploadDir, os.ModePerm)
		}

		filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
		filePath := filepath.Join(uploadDir, filename)

		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to UPDATE data",
				"errors":  []string{"Gagal menyimpan foto"},
				"data":    nil,
			})
		}

		toko.UrlFoto = filename
	}

	if err := database.DB.Save(&toko).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to UPDATE data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to UPDATE data",
		"errors":  nil,
		"data":    "Update toko succeed",
	})
}
