package models

import "time"

type Toko struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IDUser    uint      `json:"id_user"`
	NamaToko  string    `gorm:"type:varchar(255)" json:"nama_toko"`
	UrlFoto   string    `gorm:"type:varchar(255)" json:"url_foto"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi
	Products []Product `gorm:"foreignKey:IDToko"`
}
