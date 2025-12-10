package repositories

import (
	"warehouse/config"
	"warehouse/models"

	"gorm.io/gorm"
)

type StokRepository struct {
	db *gorm.DB
}

func NewStokRepository() *StokRepository {
	return &StokRepository{db: config.DBInit()}
}

func (r *StokRepository) GetAll() ([]models.MStok, error) {
	var stoks []models.MStok
	err := r.db.Preload("Barang").Find(&stoks).Error
	return stoks, err
}

func (r *StokRepository) GetByBarangID(barangID uint) (*models.MStok, error) {
	var stok models.MStok
	err := r.db.Preload("Barang").Where("barang_id = ?", barangID).First(&stok).Error
	if err != nil {
		return nil, err
	}
	return &stok, nil
}

func (r *StokRepository) CreateOrUpdate(tx *gorm.DB, barangID uint, qty int) error {
	var stok models.MStok
	err := tx.Where("barang_id = ?", barangID).First(&stok).Error

	if err == gorm.ErrRecordNotFound {
		stok = models.MStok{
			BarangID:  barangID,
			StokAkhir: qty,
		}
		return tx.Create(&stok).Error
	}

	if err != nil {
		return err
	}

	stok.StokAkhir += qty
	return tx.Save(&stok).Error
}

func (r *StokRepository) UpdateStok(tx *gorm.DB, barangID uint, newStok int) error {
	return tx.Model(&models.MStok{}).Where("barang_id = ?", barangID).Update("stok_akhir", newStok).Error
}

func (r *StokRepository) GetCurrentStok(tx *gorm.DB, barangID uint) (int, error) {
	var stok models.MStok
	err := tx.Where("barang_id = ?", barangID).First(&stok).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return stok.StokAkhir, nil
}

func (r *StokRepository) GetHistoryByBarang(barangID uint) ([]models.HistoryStok, error) {
	var history []models.HistoryStok
	err := r.db.Preload("Barang").Preload("User").Where("barang_id = ?", barangID).Order("created_at DESC").Find(&history).Error
	return history, err
}

func (r *StokRepository) GetAllHistory(page, limit int) ([]models.HistoryStok, int64, error) {
	var history []models.HistoryStok
	var total int64

	r.db.Model(&models.HistoryStok{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Preload("Barang").Preload("User").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&history).Error

	return history, total, err
}

func (r *StokRepository) CreateHistory(tx *gorm.DB, history *models.HistoryStok) error {
	return tx.Create(history).Error
}

func (r *StokRepository) EnsureStokExists(tx *gorm.DB, barangID uint) error {
	var count int64
	tx.Model(&models.MStok{}).Where("barang_id = ?", barangID).Count(&count)
	if count == 0 {
		stok := models.MStok{
			BarangID:  barangID,
			StokAkhir: 0,
		}
		return tx.Create(&stok).Error
	}
	return nil
}
