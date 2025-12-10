package repositories

import (
	"warehouse/models"

	"gorm.io/gorm"
)

type HistoryStokRepository interface {
	Create(history *models.HistoryStok) error
	GetAll() ([]models.HistoryStok, error)
	GetByBarang(barangID uint) ([]models.HistoryStok, error)
	GetByID(id uint) (models.HistoryStok, error)
}

type historyStokRepository struct {
	db *gorm.DB
}

func NewHistoryStokRepository(db *gorm.DB) HistoryStokRepository {
	return &historyStokRepository{db}
}

func (r *historyStokRepository) Create(history *models.HistoryStok) error {
	return r.db.Create(history).Error
}

func (r *historyStokRepository) GetAll() ([]models.HistoryStok, error) {
	var items []models.HistoryStok
	err := r.db.
		Preload("Barang").
		Preload("User").
		Find(&items).Error
	return items, err
}

func (r *historyStokRepository) GetByBarang(barangID uint) ([]models.HistoryStok, error) {
	var items []models.HistoryStok
	err := r.db.
		Where("barang_id = ?", barangID).
		Preload("Barang").
		Preload("User").
		Find(&items).Error
	return items, err
}

func (r *historyStokRepository) GetByID(id uint) (models.HistoryStok, error) {
	var hs models.HistoryStok
	err := r.db.
		Preload("Barang").
		Preload("User").
		First(&hs, id).Error
	return hs, err
}
