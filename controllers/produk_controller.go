package controllers

import (
	"github.com/MadeRaditya/rakamin-evermos-go/database"
	"github.com/MadeRaditya/rakamin-evermos-go/models"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func makeSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

func GetAllProduct(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	namaProduk := c.Query("nama_produk")
	categoryID := c.Query("category_id")
	tokoID := c.Query("toko_id")
	minHarga := c.Query("min_harga")
	maxHarga := c.Query("max_harga")

	var products []models.Product
	query := database.DB.Model(&models.Product{}).
		Preload("Toko").
		Preload("Category").
		Preload("Photos")

	if namaProduk != "" {
		query = query.Where("nama_produk LIKE ?", "%"+namaProduk+"%")
	}
	if categoryID != "" {
		query = query.Where("id_category = ?", categoryID)
	}
	if tokoID != "" {
		query = query.Where("id_toko = ?", tokoID)
	}
	if minHarga != "" {
		query = query.Where("CAST(harga_konsumen AS UNSIGNED) >= ?", minHarga)
	}
	if maxHarga != "" {
		query = query.Where("CAST(harga_konsumen AS UNSIGNED) <= ?", maxHarga)
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
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
			"data":  products,
			"page":  page,
			"limit": limit,
		},
	})
}

func GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := database.DB.
		Preload("Toko").
		Preload("Category").
		Preload("Photos").
		First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"No Data Product"},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    product,
	})
}

func CreateProduct(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Unauthorized"})
	}

	var toko models.Toko
	if err := database.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"Anda belum memiliki toko"},
			"data":    nil,
		})
	}

	namaProduk := c.FormValue("nama_produk")
	categoryID, _ := strconv.Atoi(c.FormValue("category_id"))
	hargaReseller := c.FormValue("harga_reseller")
	hargaKonsumen := c.FormValue("harga_konsumen")
	stok, _ := strconv.Atoi(c.FormValue("stok"))
	deskripsi := c.FormValue("deskripsi")

	if namaProduk == "" || categoryID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"Nama Produk dan Kategori wajib diisi"},
			"data":    nil,
		})
	}

	newProduct := models.Product{
		NamaProduk:    namaProduk,
		Slug:          makeSlug(namaProduk),
		HargaReseller: hargaReseller,
		HargaKonsumen: hargaKonsumen,
		Stok:          stok,
		Deskripsi:     deskripsi,
		IDToko:        toko.ID,
		IDCategory:    uint(categoryID),
	}

	if err := database.DB.Create(&newProduct).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["photos"]

		for _, file := range files {
			uploadDir := "./public/uploads/products"
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				_ = os.Mkdir(uploadDir, os.ModePerm)
			}

			filename := fmt.Sprintf("%d-%d-%s", newProduct.ID, time.Now().UnixNano(), file.Filename)
			filePath := filepath.Join(uploadDir, filename)
			c.SaveFile(file, filePath)

			photo := models.FotoProduk{
				IDProduk: newProduct.ID,
				Url:      filename,
			}
			database.DB.Create(&photo)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data":    newProduct.ID,
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")

	var toko models.Toko
	if err := database.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false, "message": "Failed to PUT data", "errors": []string{"Toko tidak ditemukan"}, "data": nil,
		})
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND id_toko = ?", id, toko.ID).First(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"record not found"},
			"data":    nil,
		})
	}

	if val := c.FormValue("nama_produk"); val != "" {
		product.NamaProduk = val
		product.Slug = makeSlug(val)
	}
	if val := c.FormValue("harga_reseller"); val != "" {
		product.HargaReseller = val
	}
	if val := c.FormValue("harga_konsumen"); val != "" {
		product.HargaKonsumen = val
	}
	if val := c.FormValue("deskripsi"); val != "" {
		product.Deskripsi = val
	}
	if val := c.FormValue("stok"); val != "" {
		stok, _ := strconv.Atoi(val)
		product.Stok = stok
	}
	if val := c.FormValue("category_id"); val != "" {
		catID, _ := strconv.Atoi(val)
		product.IDCategory = uint(catID)
	}

	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false, "message": "Failed to UPDATE data", "errors": []string{err.Error()}, "data": nil,
		})
	}

	form, err := c.MultipartForm()
	if err == nil {
		files := form.File["photos"]
		if len(files) > 0 {
			for _, file := range files {
				uploadDir := "./public/uploads/products"
				if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
					_ = os.Mkdir(uploadDir, os.ModePerm)
				}

				filename := fmt.Sprintf("%d-%d-%s", product.ID, time.Now().UnixNano(), file.Filename)
				filePath := filepath.Join(uploadDir, filename)
				c.SaveFile(file, filePath)

				photo := models.FotoProduk{IDProduk: product.ID, Url: filename}
				database.DB.Create(&photo)
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    "",
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")

	var toko models.Toko
	if err := database.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false, "message": "Failed", "errors": []string{"Unauthorized"}, "data": nil,
		})
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND id_toko = ?", id, toko.ID).First(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"record not found"},
			"data":    nil,
		})
	}

	database.DB.Where("id_produk = ?", product.ID).Delete(&models.FotoProduk{})

	if err := database.DB.Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false, "message": "Failed to DELETE", "errors": []string{err.Error()}, "data": nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    "",
	})
}
