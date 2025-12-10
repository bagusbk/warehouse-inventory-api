package repositories

import (
	"fmt"
	"time"
	"warehouse/config"
	"warehouse/models"

	"gorm.io/gorm"
)

type PembelianRepository struct {
	db *gorm.DB
}

func NewPembelianRepository() *PembelianRepository {
	return &PembelianRepository{db: config.DBInit()}
}

func (r *PembelianRepository) Create(req *models.PembelianRequest, userID uint) (*models.BeliHeader, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	noFaktur := r.generateNoFaktur()

	var total float64
	for _, detail := range req.Details {
		total += float64(detail.Qty) * detail.Harga
	}

	header := models.BeliHeader{
		NoFaktur: noFaktur,
		Supplier: req.Supplier,
		Total:    total,
		UserID:   userID,
		Status:   "selesai",
	}

	if err := tx.Create(&header).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	stokRepo := NewStokRepository()

	for _, detailReq := range req.Details {
		var barang models.MasterBarang
		if err := tx.First(&barang, detailReq.BarangID).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("barang dengan ID %d tidak ditemukan", detailReq.BarangID)
		}

		subtotal := float64(detailReq.Qty) * detailReq.Harga
		detail := models.BeliDetail{
			BeliHeaderID: header.ID,
			BarangID:     detailReq.BarangID,
			Qty:          detailReq.Qty,
			Harga:        detailReq.Harga,
			Subtotal:     subtotal,
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := stokRepo.EnsureStokExists(tx, detailReq.BarangID); err != nil {
			tx.Rollback()
			return nil, err
		}

		stokSebelum, err := stokRepo.GetCurrentStok(tx, detailReq.BarangID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		stokSesudah := stokSebelum + detailReq.Qty

		if err := stokRepo.UpdateStok(tx, detailReq.BarangID, stokSesudah); err != nil {
			tx.Rollback()
			return nil, err
		}

		history := models.HistoryStok{
			BarangID:       detailReq.BarangID,
			UserID:         userID,
			JenisTransaksi: "masuk",
			Jumlah:         detailReq.Qty,
			StokSebelum:    stokSebelum,
			StokSesudah:    stokSesudah,
			Keterangan:     fmt.Sprintf("Pembelian %s", noFaktur),
		}

		if err := stokRepo.CreateHistory(tx, &history); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &header, nil
}

func (r *PembelianRepository) GetAll(page, limit int, startDate, endDate string) ([]models.BeliHeader, int64, error) {
	var headers []models.BeliHeader
	var total int64

	query := r.db.Model(&models.BeliHeader{})

	if startDate != "" && endDate != "" {
		query = query.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("User").Preload("Details").Preload("Details.Barang").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&headers).Error

	return headers, total, err
}

func (r *PembelianRepository) GetByID(id uint) (*models.BeliHeader, error) {
	var header models.BeliHeader
	err := r.db.Preload("User").Preload("Details").Preload("Details.Barang").First(&header, id).Error
	if err != nil {
		return nil, err
	}
	return &header, nil
}

func (r *PembelianRepository) generateNoFaktur() string {
	now := time.Now()
	var count int64
	r.db.Model(&models.BeliHeader{}).
		Where("DATE(created_at) = ?", now.Format("2006-01-02")).
		Count(&count)

	return fmt.Sprintf("BLI%s%03d", now.Format("20060102"), count+1)
}
