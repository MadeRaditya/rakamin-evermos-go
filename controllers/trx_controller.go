package controllers

import (
	"backend-evermos/database"
	"backend-evermos/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DetailTrxInput struct {
	ProductID uint `json:"product_id"`
	Kuantitas int  `json:"kuantitas"`
}

type TrxInput struct {
	MethodBayar string           `json:"method_bayar"`
	AlamatKirim uint             `json:"alamat_kirim"`
	DetailTrx   []DetailTrxInput `json:"detail_trx"`
}

func GetAllTrx(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Unauthorized"})
	}

	var trxs []models.Trx

	if err := database.DB.
		Preload("Alamat").
		Preload("DetailTrx").
		Preload("DetailTrx.Toko").
		Preload("DetailTrx.Product").
		Preload("DetailTrx.Product.Toko").
		Preload("DetailTrx.Product.Category").
		Preload("DetailTrx.Product.Photos").
		Where("id_user = ?", userID).
		Find(&trxs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false, "message": "Failed to GET data", "errors": []string{err.Error()}, "data": nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data": fiber.Map{
			"data": trxs,
		},
	})
}

func GetTrxByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id")

	var trxID uint
	if _, err := fmt.Sscan(id, &trxID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "ID Invalid"})
	}

	var trx models.Trx
	if err := database.DB.
		Preload("Alamat").
		Preload("DetailTrx").
		Preload("DetailTrx.Toko").
		Preload("DetailTrx.Product").
		Preload("DetailTrx.Product.Toko").
		Preload("DetailTrx.Product.Category").
		Preload("DetailTrx.Product.Photos").
		Where("id = ? AND id_user = ?", trxID, userID).
		First(&trx).Error; err != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  []string{"No Data Trx"},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data":    trx,
	})
}

func CreateTrx(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "Unauthorized"})
	}
	uid := uint(userID.(float64))

	var input TrxInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Failed to POST data", "errors": []string{"Input invalid"}, "data": nil,
		})
	}

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	totalHargaTrx := 0
	var newDetailTrx []models.DetailTrx

	for _, item := range input.DetailTrx {
		var product models.Product

		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": false, "message": "Failed", "errors": []string{"Product not found"}, "data": nil,
			})
		}

		if product.Stok < item.Kuantitas {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": false, "message": "Failed",
				"errors": []string{"Stok tidak mencukupi untuk produk: " + product.NamaProduk},
				"data":   nil,
			})
		}

		product.Stok -= item.Kuantitas
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "message": "Failed stock update"})
		}

		logProduk := models.LogProduk{
			IDProduk:      product.ID,
			NamaProduk:    product.NamaProduk,
			Slug:          product.Slug,
			HargaReseller: product.HargaReseller,
			HargaKonsumen: product.HargaKonsumen,
			Deskripsi:     product.Deskripsi,
			IDToko:        product.IDToko,
			IDCategory:    product.IDCategory,
		}

		if err := tx.Create(&logProduk).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "message": "Failed create log"})
		}

		hargaKonsumenInt, _ := strconv.Atoi(product.HargaKonsumen)
		subTotal := hargaKonsumenInt * item.Kuantitas
		totalHargaTrx += subTotal

		newDetailTrx = append(newDetailTrx, models.DetailTrx{
			IDLogProduk: logProduk.ID,
			IDToko:      product.IDToko,
			Kuantitas:   item.Kuantitas,
			HargaTotal:  subTotal,
		})
	}

	newTrx := models.Trx{
		IDUser:           uid,
		AlamatPengiriman: input.AlamatKirim,
		HargaTotal:       totalHargaTrx,
		KodeInvoice:      fmt.Sprint("INV-%d", time.Now().Unix()),
		MethodBayar:      input.MethodBayar,
	}

	if err := tx.Create(&newTrx).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "message": "Failed create trx"})
	}

	for i := range newDetailTrx {
		newDetailTrx[i].IDTrx = newTrx.ID
		if err := tx.Create(&newDetailTrx[i]).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "message": "Failed create detail"})
		}
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data":    newTrx.ID,
	})
}
