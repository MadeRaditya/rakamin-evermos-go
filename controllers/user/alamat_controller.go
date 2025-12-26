package user

import (
	"github.com/MadeRaditya/rakamin-evermos-go/database"
	"github.com/MadeRaditya/rakamin-evermos-go/models"
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

func GetMyAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"Token tidak valid"},
			"data":    nil,
		})
	}

	var alamat []models.Alamat
	query := database.DB.Where("id_user = ?", userID)

	if judul := c.Query("judul_alamat"); judul != "" {
		query = query.Where("judul_alamat LIKE ?", "%"+judul+"%")
	}

	if err := query.Find(&alamat).Error; err != nil {
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
		"data":    alamat,
	})
}

func CreateAlamat(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"Token tidak valid"},
			"data":    nil,
		})
	}

	var input AlamatInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"Input tidak valid"},
			"data":    nil,
		})
	}

	uidFloat, ok := userID.(float64)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false, "message": "Failed to POST data", "errors": []string{"User ID Error"}, "data": nil,
		})
	}
	uid := uint(uidFloat)

	alamat := models.Alamat{
		IDUser:       uid,
		JudulAlamat:  input.JudulAlamat,
		NamaPenerima: input.NamaPenerima,
		NoTelp:       input.NoTelp,
		DetailAlamat: input.DetailAlamat,
	}

	if err := database.DB.Create(&alamat).Error; err != nil {
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
		"data":    alamat.ID,
	})
}

func GetAlamatByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")

	var alamatID uint
	if _, err := fmt.Sscan(id, &alamatID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"ID tidak valid"},
			"data":    nil,
		})
	}

	var alamat models.Alamat
	if err := database.DB.Where("id = ? AND id_user = ?", alamatID, userID).First(&alamat).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"record not found"},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    alamat,
	})
}

func UpdateAlamat(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")
	var input AlamatUpdate

	var alamatID uint
	if _, err := fmt.Sscan(id, &alamatID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"ID tidak valid"},
			"data":    nil,
		})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"Input tidak valid"},
			"data":    nil,
		})
	}

	var alamat models.Alamat
	if err := database.DB.Where("id = ? AND id_user = ?", alamatID, userID).First(&alamat).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"record not found"},
			"data":    nil,
		})
	}

	alamat.NamaPenerima = input.NamaPenerima
	alamat.NoTelp = input.NoTelp
	alamat.DetailAlamat = input.DetailAlamat

	if err := database.DB.Save(&alamat).Error; err != nil {
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

func DeleteAlamat(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")

	var alamatID uint
	if _, err := fmt.Sscan(id, &alamatID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"ID tidak valid"},
			"data":    nil,
		})
	}

	var alamat models.Alamat
	if err := database.DB.Where("id = ? AND id_user = ?", alamatID, userID).First(&alamat).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"record not found"},
			"data":    nil,
		})
	}

	if err := database.DB.Delete(&alamat).Error; err != nil {
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
