package database

import (
	"github.com/MadeRaditya/rakamin-evermos-go/models"
	"fmt"

	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_CHARSET"),
		os.Getenv("DB_PARSE_TIME"),
		os.Getenv("DB_LOC"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal koneksi ke database")
	}

	db.AutoMigrate(
		&models.User{},
		&models.Alamat{},
		&models.Toko{},
		&models.Category{},
		&models.Product{},
		&models.FotoProduk{},
		&models.Trx{},
		&models.DetailTrx{},
		&models.LogProduk{},
	)

	DB = db
	fmt.Println("Database Connected & Migrated")
}
