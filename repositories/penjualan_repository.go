package repositories

import (
	"errors"
	"fmt"
	"time"
	"warehouse/config"
	"warehouse/models"

	"gorm.io/gorm"
)

type PenjualanRepository struct {
	db *gorm.DB
}

func NewPenjualanRepository() *PenjualanRepository {
	return &PenjualanRepository{db: config.DBInit()}
}

func (r *PenjualanRepository) Create(req *models.PenjualanRequest, userID uint) (*models.JualHeader, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	stokRepo := NewStokRepository()

	for _, detailReq := range req.Details {
		var barang models.MasterBarang
		if err := tx.First(&barang, detailReq.BarangID).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("barang dengan ID %d tidak ditemukan", detailReq.BarangID)
		}

		currentStok, err := stokRepo.GetCurrentStok(tx, detailReq.BarangID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if currentStok < detailReq.Qty {
			tx.Rollback()
			return nil, errors.New(fmt.Sprintf("stok tidak mencukupi untuk barang %s. Stok tersedia: %d, diminta: %d",
				barang.NamaBarang, currentStok, detailReq.Qty))
		}
	}

	noFaktur := r.generateNoFaktur()

	var total float64
	for _, detail := range req.Details {
		total += float64(detail.Qty) * detail.Harga
	}

	header := models.JualHeader{
		NoFaktur: noFaktur,
		Customer: req.Customer,
		Total:    total,
		UserID:   userID,
		Status:   "selesai",
	}

	if err := tx.Create(&header).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, detailReq := range req.Details {
		subtotal := float64(detailReq.Qty) * detailReq.Harga
		detail := models.JualDetail{
			JualHeaderID: header.ID,
			BarangID:     detailReq.BarangID,
			Qty:          detailReq.Qty,
			Harga:        detailReq.Harga,
			Subtotal:     subtotal,
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		stokSebelum, _ := stokRepo.GetCurrentStok(tx, detailReq.BarangID)
		stokSesudah := stokSebelum - detailReq.Qty

		if err := stokRepo.UpdateStok(tx, detailReq.BarangID, stokSesudah); err != nil {
			tx.Rollback()
			return nil, err
		}

		history := models.HistoryStok{
			BarangID:       detailReq.BarangID,
			UserID:         userID,
			JenisTransaksi: "keluar",
			Jumlah:         detailReq.Qty,
			StokSebelum:    stokSebelum,
			StokSesudah:    stokSesudah,
			Keterangan:     fmt.Sprintf("Penjualan %s", noFaktur),
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

func (r *PenjualanRepository) GetAll(page, limit int, startDate, endDate string) ([]models.JualHeader, int64, error) {
	var headers []models.JualHeader
	var total int64

	query := r.db.Model(&models.JualHeader{})

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

func (r *PenjualanRepository) GetByID(id uint) (*models.JualHeader, error) {
	var header models.JualHeader
	err := r.db.Preload("User").Preload("Details").Preload("Details.Barang").First(&header, id).Error
	if err != nil {
		return nil, err
	}
	return &header, nil
}

func (r *PenjualanRepository) generateNoFaktur() string {
	now := time.Now()
	var count int64
	r.db.Model(&models.JualHeader{}).
		Where("DATE(created_at) = ?", now.Format("2006-01-02")).
		Count(&count)

	return fmt.Sprintf("JUAL%s%03d", now.Format("20060102"), count+1)
}
