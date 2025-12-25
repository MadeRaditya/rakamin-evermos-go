package models

import "time"

type Alamat struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	IDUser       uint      `json:"id_user"`
	JudulAlamat  string    `gorm:"type:varchar(255)" json:"judul_alamat"`
	NamaPenerima string    `gorm:"type:varchar(255)" json:"nama_penerima"`
	NoTelp       string    `gorm:"type:varchar(255)" json:"no_telp"`
	DetailAlamat string    `gorm:"type:varchar(255)" json:"detail_alamat"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
