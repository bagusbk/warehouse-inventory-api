package models

import "time"

type HistoryStok struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	BarangID       uint         `json:"barang_id"`
	UserID         uint         `json:"user_id"`
	JenisTransaksi string       `json:"jenis_transaksi" gorm:"not null;size:50"`
	Jumlah         int          `json:"jumlah" gorm:"not null"`
	StokSebelum    int          `json:"stok_sebelum" gorm:"not null"`
	StokSesudah    int          `json:"stok_sesudah" gorm:"not null"`
	Keterangan     string       `json:"keterangan" gorm:"type:text"`
	CreatedAt      time.Time    `json:"created_at"`
	Barang         MasterBarang `json:"barang,omitempty" gorm:"foreignKey:BarangID"`
	User           User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (HistoryStok) TableName() string {
	return "history_stok"
}

type HistoryStokResponse struct {
	ID             uint         `json:"id"`
	BarangID       uint         `json:"barang_id"`
	UserID         uint         `json:"user_id"`
	JenisTransaksi string       `json:"jenis_transaksi"`
	Jumlah         int          `json:"jumlah"`
	StokSebelum    int          `json:"stok_sebelum"`
	StokSesudah    int          `json:"stok_sesudah"`
	Keterangan     string       `json:"keterangan"`
	CreatedAt      time.Time    `json:"created_at"`
	Barang         BarangSimple `json:"barang"`
	User           UserResponse `json:"user"`
}
