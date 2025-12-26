package controllers

import (
	"github.com/MadeRaditya/rakamin-evermos-go/config"
	"github.com/MadeRaditya/rakamin-evermos-go/database"
	"github.com/MadeRaditya/rakamin-evermos-go/models"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false,
			"message": "Failed to POST data",
			"errors":  []string{"Input tidak valid"},
			"data":    nil,
		})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.KataSandi), bcrypt.DefaultCost)

	var tglLahir time.Time
	if strings.TrimSpace(input.TanggalLahir) != "" {
		parsedDate, err := time.Parse("02/01/2006", input.TanggalLahir)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to POST data",
				"errors":  []string{"Format tanggal_Lahir harus DD/MM/YYYY"},
				"data":    nil,
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
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	toko := models.Toko{
		IDUser:   user.ID,
		NamaToko: "Toko " + user.Nama,
		UrlFoto:  "",
	}
	database.DB.Create(&toko)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data":    "Register Succeed",
	})
}

func Login(c *fiber.Ctx) error {
	var input LoginInput
	var user models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"Input tidak valid"},
			"data":    nil,
		})
	}

	if err := database.DB.Where("no_telp = ?", input.NoTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"No Telp atau Password salah"},
			"data":    nil,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(input.KataSandi)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"No Telp atau Password salah"},
			"data":    nil,
		})
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(config.JWTSecret())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to generate token",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	responseData := fiber.Map{
		"nama":          user.Nama,
		"no_telp":       user.NoTelp,
		"tanggal_Lahir": user.TanggalLahir.Format("02/01/2006"),
		"tentang":       user.Tentang,
		"pekerjaan":     user.Pekerjaan,
		"email":         user.Email,
		"id_provinsi":   user.IDProvinsi,
		"id_kota":       user.IDKota,
		"token":         t,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data":    responseData,
	})
}
