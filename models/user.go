package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Nama         string    `gorm:"type:varchar(255)" json:"nama"`
	KataSandi    string    `gorm:"type:varchar(255)" json:"-"`
	NoTelp       string    `gorm:"type:varchar(255);unique" json:"no_telp"`
	TanggalLahir time.Time `gorm:"type:date" json:"tanggal_lahir"`
	JenisKelamin string    `gorm:"type:varchar(255)" json:"jenis_kelamin"`
	Tentang      string    `gorm:"type:text" json:"tentang"`
	Pekerjaan    string    `gorm:"type:varchar(255)" json:"pekerjaan"`
	Email        string    `gorm:"type:varchar(255);unique" json:"email"`
	IDProvinsi   string    `gorm:"type:varchar(255)" json:"id_provinsi"`
	IDKota       string    `gorm:"type:varchar(255)" json:"id_kota"`
	IsAdmin      bool      `gorm:"default:false" json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relasi
	Alamat []Alamat `gorm:"foreignKey:IDUser"` // One to Many
	Toko   Toko     `gorm:"foreignKey:IDUser"` // One to One
}
