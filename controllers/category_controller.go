package controllers

import (
	"github.com/MadeRaditya/rakamin-evermos-go/database"
	"github.com/MadeRaditya/rakamin-evermos-go/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type CategoryInput struct {
	NamaCategory string `json:"nama_category"`
}

func isAdmin(c *fiber.Ctx) bool {
	userIDLocals := c.Locals("user_id")
	if userIDLocals == nil {
		return false
	}

	userIDFloat, ok := userIDLocals.(float64)
	if !ok {
		return false
	}

	var user models.User
	if err := database.DB.First(&user, uint(userIDFloat)).Error; err != nil {
		return false
	}
	return user.IsAdmin
}

func GetAllCategory(c *fiber.Ctx) error {
	var categories []models.Category

	if err := database.DB.Find(&categories).Error; err != nil {
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
		"data":    categories,
	})
}

func CreateCategory(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"Unauthorized"},
			"data":    nil,
		})
	}

	var categoryInput CategoryInput

	if err := c.BodyParser(&categoryInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"Input tidak valid"},
			"data":    nil,
		})
	}

	if categoryInput.NamaCategory == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"Nama category wajib diisi"},
			"data":    nil,
		})
	}

	category := models.Category{
		NamaCategory: categoryInput.NamaCategory,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data":    category.ID,
	})
}

func GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var categoryID uint
	if _, err := fmt.Sscan(id, &categoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"ID tidak valid"},
			"data":    nil,
		})
	}

	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"No Data Category"},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    category,
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{"Unauthorized"},
			"data":    nil,
		})
	}

	id := c.Params("id")
	var categoryID uint
	if _, err := fmt.Sscan(id, &categoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{"ID tidak valid"},
			"data":    nil,
		})
	}

	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"record not found"},
			"data":    nil,
		})
	}

	var updateData CategoryInput
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to PUT data",
			"errors":  []string{"Invalid request body"},
			"data":    nil,
		})
	}

	if updateData.NamaCategory != "" {
		category.NamaCategory = updateData.NamaCategory
	}

	if err := database.DB.Save(&category).Error; err != nil {
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

func DeleteCategory(c *fiber.Ctx) error {
	if !isAdmin(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to DELETE data",
			"errors":  []string{"Unauthorized"},
			"data":    nil,
		})
	}

	id := c.Params("id")
	var categoryID uint
	if _, err := fmt.Sscan(id, &categoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to DELETE data",
			"errors":  []string{"ID tidak valid"},
			"data":    nil,
		})
	}

	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"record not found"},
			"data":    nil,
		})
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to DELETE data",
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
