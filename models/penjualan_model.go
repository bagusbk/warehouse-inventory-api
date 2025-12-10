package models

import "time"

type JualHeader struct {
	ID        uint         `json:"id" gorm:"primaryKey"`
	NoFaktur  string       `json:"no_faktur" gorm:"unique;not null;size:100"`
	Customer  string       `json:"customer" gorm:"not null;size:200"`
	Total     float64      `json:"total" gorm:"type:decimal(15,2);default:0"`
	UserID    uint         `json:"user_id"`
	Status    string       `json:"status" gorm:"default:selesai;size:50"`
	CreatedAt time.Time    `json:"created_at"`
	User      User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Details   []JualDetail `json:"details,omitempty" gorm:"foreignKey:JualHeaderID"`
}

func (JualHeader) TableName() string {
	return "jual_header"
}

type JualDetail struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	JualHeaderID uint         `json:"jual_header_id"`
	BarangID     uint         `json:"barang_id"`
	Qty          int          `json:"qty" gorm:"not null"`
	Harga        float64      `json:"harga" gorm:"type:decimal(15,2);not null"`
	Subtotal     float64      `json:"subtotal" gorm:"type:decimal(15,2);not null"`
	Barang       MasterBarang `json:"barang,omitempty" gorm:"foreignKey:BarangID"`
}

func (JualDetail) TableName() string {
	return "jual_detail"
}

type PenjualanRequest struct {
	Customer string                   `json:"customer" binding:"required"`
	Details  []PenjualanDetailRequest `json:"details" binding:"required,dive"`
}

type PenjualanDetailRequest struct {
	BarangID uint    `json:"barang_id" binding:"required"`
	Qty      int     `json:"qty" binding:"required,min=1"`
	Harga    float64 `json:"harga" binding:"required"`
}

type PenjualanResponse struct {
	Header  JualHeaderResponse   `json:"header"`
	Details []JualDetailResponse `json:"details"`
}

type JualHeaderResponse struct {
	ID        uint         `json:"id"`
	NoFaktur  string       `json:"no_faktur"`
	Customer  string       `json:"customer"`
	Total     float64      `json:"total"`
	UserID    uint         `json:"user_id"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	User      UserResponse `json:"user"`
}

type JualDetailResponse struct {
	ID       uint         `json:"id"`
	BarangID uint         `json:"barang_id"`
	Qty      int          `json:"qty"`
	Harga    float64      `json:"harga"`
	Subtotal float64      `json:"subtotal"`
	Barang   BarangSimple `json:"barang"`
}
