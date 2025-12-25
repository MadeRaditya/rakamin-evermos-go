package controllers

import (
	"backend-evermos/database"
	"backend-evermos/models"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CategoryInput struct {
	NamaCategory string `json:"nama_category"`
}

func GetAllCategory(c *fiber.Ctx) error {
	var categories []models.Category

	if err := database.DB.Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": categories,
	})
}

func CreateCategory(c *fiber.Ctx) error {
	var categoryInput CategoryInput

	if err := c.BodyParser(&categoryInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	if categoryInput.NamaCategory == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nama category wajib diisi",
		})
	}

	category := models.Category{
		NamaCategory: categoryInput.NamaCategory,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambahkan category.",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category berhasil ditambahkan",
		"data": fiber.Map{
			"id":            category.ID,
			"nama_category": category.NamaCategory,
		},
	})
}

func GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var categoryID uint
	if _, err := fmt.Sscan(id, &categoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
		})
	}

	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Category tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": category,
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var categoryID uint
	if _, err := fmt.Sscan(id, &categoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
		})
	}

	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Category tidak ditemukan",
		})
	}

	var updateData CategoryInput
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if strings.TrimSpace(updateData.NamaCategory) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nama category wajib diisi",
		})
	}

	category.NamaCategory = updateData.NamaCategory

	if err := database.DB.Model(&category).Updates(models.Category{
		NamaCategory: updateData.NamaCategory}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Message": "Gagal memperbarui category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category berhasil di perbarui",
		"data":    category,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var categoryID uint
	if _, err := fmt.Sscan(id, &categoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID tidak valid",
		})
	}

	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Category tidak ditemukan",
		})
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category berhasil dihapus",
		"data":    category,
	})
}
