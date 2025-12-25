package controllers

import (
	"backend-evermos/config"
	"backend-evermos/database"
	"backend-evermos/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Nama         string `json:"nama"`
	NoTelp       string `json:"no_telp"`
	KataSandi    string `json:"kata_sandi"`
	TanggalLahir string `json:"tanggal_Lahir"`
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}

type LoginInput struct {
	NoTelp    string `json:"no_telp"`
	KataSandi string `json:"kata_sandi"`
}

func Register(c *fiber.Ctx) error {
	var input RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.KataSandi), bcrypt.DefaultCost)

	var tglLahir time.Time
	if strings.TrimSpace(input.TanggalLahir) != "" {
		parsedDate, err := time.Parse("02/01/2006", input.TanggalLahir)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Format tanggal_Lahir harus DD/MM/YYYY",
			})
		}
		tglLahir = parsedDate
	}

	user := models.User{
		Nama:         input.Nama,
		NoTelp:       input.NoTelp,
		KataSandi:    string(hashedPassword),
		TanggalLahir: tglLahir,
		Pekerjaan:    input.Pekerjaan,
		Email:        input.Email,
		IDProvinsi:   input.IDProvinsi,
		IDKota:       input.IDKota,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Gagal registrasi. Email atau No Telp mungkin sudah terdaftar.",
		})
	}

	toko := models.Toko{
		IDUser:   user.ID,
		NamaToko: "Toko " + user.Nama,
		UrlFoto:  "",
	}
	database.DB.Create(&toko)

	var result models.User
	database.DB.Preload("Toko").First(&result, user.ID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Register Berhasil",
		"data":    result,
	})
}

func Login(c *fiber.Ctx) error {
	var input LoginInput
	var user models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input tidak valid"})
	}

	if err := database.DB.Where("no_telp = ?", input.NoTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "No Telp atau Password salah"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(input.KataSandi)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "No Telp atau Password salah"})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal generate token"})
	}

	var result models.User
	database.DB.Preload("Toko").First(&result, user.ID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"token": t,
			"user":  result,
		},
		"message": "Login Berhasil",
	})
}
