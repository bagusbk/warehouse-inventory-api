package models

import "time"

type BeliHeader struct {
	ID        uint         `json:"id" gorm:"primaryKey"`
	NoFaktur  string       `json:"no_faktur" gorm:"unique;not null;size:100"`
	Supplier  string       `json:"supplier" gorm:"not null;size:200"`
	Total     float64      `json:"total" gorm:"type:decimal(15,2);default:0"`
	UserID    uint         `json:"user_id"`
	Status    string       `json:"status" gorm:"default:selesai;size:50"`
	CreatedAt time.Time    `json:"created_at"`
	User      User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Details   []BeliDetail `json:"details,omitempty" gorm:"foreignKey:BeliHeaderID"`
}

func (BeliHeader) TableName() string {
	return "beli_header"
}

type BeliDetail struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	BeliHeaderID uint         `json:"beli_header_id"`
	BarangID     uint         `json:"barang_id"`
	Qty          int          `json:"qty" gorm:"not null"`
	Harga        float64      `json:"harga" gorm:"type:decimal(15,2);not null"`
	Subtotal     float64      `json:"subtotal" gorm:"type:decimal(15,2);not null"`
	Barang       MasterBarang `json:"barang,omitempty" gorm:"foreignKey:BarangID"`
}

func (BeliDetail) TableName() string {
	return "beli_detail"
}

type PembelianRequest struct {
	Supplier string                   `json:"supplier" binding:"required"`
	Details  []PembelianDetailRequest `json:"details" binding:"required,dive"`
}

type PembelianDetailRequest struct {
	BarangID uint    `json:"barang_id" binding:"required"`
	Qty      int     `json:"qty" binding:"required,min=1"`
	Harga    float64 `json:"harga" binding:"required"`
}

type PembelianResponse struct {
	Header  BeliHeaderResponse   `json:"header"`
	Details []BeliDetailResponse `json:"details"`
}

type BeliHeaderResponse struct {
	ID        uint         `json:"id"`
	NoFaktur  string       `json:"no_faktur"`
	Supplier  string       `json:"supplier"`
	Total     float64      `json:"total"`
	UserID    uint         `json:"user_id"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	User      UserResponse `json:"user"`
}

type BeliDetailResponse struct {
	ID       uint         `json:"id"`
	BarangID uint         `json:"barang_id"`
	Qty      int          `json:"qty"`
	Harga    float64      `json:"harga"`
	Subtotal float64      `json:"subtotal"`
	Barang   BarangSimple `json:"barang"`
}
