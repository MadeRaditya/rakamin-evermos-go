package models

import "time"

type LogProduk struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	IDProduk      uint      `json:"id_produk"`
	NamaProduk    string    `gorm:"type:varchar(255)" json:"nama_produk"`
	Slug          string    `gorm:"type:varchar(255)" json:"slug"`
	HargaReseller string    `gorm:"type:varchar(255)" json:"harga_reseller"`
	HargaKonsumen string    `gorm:"type:varchar(255)" json:"harga_konsumen"`
	Deskripsi     string    `gorm:"type:text" json:"deskripsi"`
	IDToko        uint      `json:"id_toko"`
	IDCategory    uint      `json:"id_category"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Toko     Toko         `gorm:"foreignKey:IDToko" json:"toko"`
	Category Category     `gorm:"foreignKey:IDCategory" json:"category"`
	Photos   []FotoProduk `gorm:"foreignKey:IDProduk;references:IDProduk" json:"photos"`
}
