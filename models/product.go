package models

import "time"

type Product struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NamaProduk    string    `gorm:"type:varchar(255)" json:"nama_produk"`
	Slug          string    `gorm:"type:varchar(255)" json:"slug"`
	HargaReseller string    `gorm:"type:varchar(255)" json:"harga_reseller"`
	HargaKonsumen string    `gorm:"type:varchar(255)" json:"harga_konsumen"`
	Stok          int       `json:"stok"`
	Deskripsi     string    `gorm:"type:text" json:"deskripsi"`
	IDToko        uint      `json:"id_toko"`
	IDCategory    uint      `json:"id_category"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relasi
	Photos []FotoProduk `gorm:"foreignKey:IDProduk"`
}
