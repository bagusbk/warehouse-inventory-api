package models

import "time"

type MasterBarang struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	KodeBarang string    `json:"kode_barang" gorm:"unique;not null;size:50"`
	NamaBarang string    `json:"nama_barang" gorm:"not null;size:200"`
	Deskripsi  string    `json:"deskripsi" gorm:"type:text"`
	Satuan     string    `json:"satuan" gorm:"not null;size:50"`
	HargaBeli  float64   `json:"harga_beli" gorm:"type:decimal(15,2);default:0"`
	HargaJual  float64   `json:"harga_jual" gorm:"type:decimal(15,2);default:0"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (MasterBarang) TableName() string {
	return "master_barang"
}

type BarangRequest struct {
	KodeBarang string  `json:"kode_barang" binding:"required"`
	NamaBarang string  `json:"nama_barang" binding:"required"`
	Deskripsi  string  `json:"deskripsi"`
	Satuan     string  `json:"satuan" binding:"required"`
	HargaBeli  float64 `json:"harga_beli"`
	HargaJual  float64 `json:"harga_jual"`
}

type BarangWithStok struct {
	ID         uint    `json:"id"`
	KodeBarang string  `json:"kode_barang"`
	NamaBarang string  `json:"nama_barang"`
	Deskripsi  string  `json:"deskripsi"`
	Satuan     string  `json:"satuan"`
	HargaBeli  float64 `json:"harga_beli"`
	HargaJual  float64 `json:"harga_jual"`
	StokAkhir  int     `json:"stok_akhir"`
}

type BarangSimple struct {
	KodeBarang string  `json:"kode_barang"`
	NamaBarang string  `json:"nama_barang"`
	Satuan     string  `json:"satuan,omitempty"`
	HargaJual  float64 `json:"harga_jual,omitempty"`
}
