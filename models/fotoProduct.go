package models

import "time"

type FotoProduk struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IDProduk  uint      `json:"id_produk"`
	Url       string    `gorm:"type:varchar(255)" json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
