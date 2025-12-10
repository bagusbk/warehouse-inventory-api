package models

import "time"

type MStok struct {
	ID        uint         `json:"id" gorm:"primaryKey"`
	BarangID  uint         `json:"barang_id"`
	StokAkhir int          `json:"stok_akhir" gorm:"default:0"`
	UpdatedAt time.Time    `json:"updated_at"`
	Barang    MasterBarang `json:"barang,omitempty" gorm:"foreignKey:BarangID"`
}

func (MStok) TableName() string {
	return "mstok"
}

type StokResponse struct {
	ID        uint         `json:"id"`
	BarangID  uint         `json:"barang_id"`
	StokAkhir int          `json:"stok_akhir"`
	UpdatedAt time.Time    `json:"updated_at"`
	Barang    BarangSimple `json:"barang"`
}
