package controllers

import (
	"backend-evermos/database"
	"backend-evermos/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductInput struct {
	NamaProduk    string `json:"nama_produk"`
	IDCategory    uint   `json:"id_category"`
	HargaReseller string `json:"harga_reseller"`
	HargaKonsumen string `json:"harga_konsumen"`
	Stok          int    `json:"stok"`
	Deskripsi     string `json:"deskripsi"`
}

func GetAllProduk(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	// nama_produk := c.Query("nama")
	// category_id := strconv.Atoi(c.Query("category_id", "2"))
	// toko_id := strconv.Atoi(c.Query("toko_id", "4"))
	// max_harga := strconv.Atoi(c.Query("max_harga", "100000"))
	// min_harga := strconv.Atoi(c.Query("min_harga", "1500"))

	offset := (page - 1) * limit

	var produk []models.Product
	query := database.DB.Model(&models.Product{})

	// if nama != "" {
	// 	query = query.Where("nama_toko LIKE ?", "%"+nama+"%")
	// }
	if err := query.
		Limit(limit).
		Offset(offset).
		Find(&produk).Error; err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data produk",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  produk,
		"page":  page,
		"limit": limit,
	})
}

func CreateProduct(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid",
		})
	}

	var input ProductInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input tidak valid",
		})
	}

	produk := models.Product{
		NamaProduk:    input.NamaProduk,
		IDCategory:    input.IDCategory,
		HargaReseller: input.HargaReseller,
		HargaKonsumen: input.HargaKonsumen,
		Stok:          input.Stok,
		Deskripsi:     input.Deskripsi,
	}

	if err := database.DB.Create(&produk).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambahkan Produk",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Produk berhasil ditambahkan",
		"data":    produk,
	})
}
