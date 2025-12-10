package repositories

import (
	"warehouse/config"
	"warehouse/models"

	"gorm.io/gorm"
)

type BarangRepository struct {
	db *gorm.DB
}

func NewBarangRepository() *BarangRepository {
	return &BarangRepository{db: config.DBInit()}
}

func (r *BarangRepository) Create(barang *models.MasterBarang) error {
	return r.db.Create(barang).Error
}

func (r *BarangRepository) Update(barang *models.MasterBarang) error {
	return r.db.Save(barang).Error
}

func (r *BarangRepository) Delete(id uint) error {
	return r.db.Delete(&models.MasterBarang{}, id).Error
}

func (r *BarangRepository) FindByID(id uint) (*models.MasterBarang, error) {
	var barang models.MasterBarang
	err := r.db.First(&barang, id).Error
	if err != nil {
		return nil, err
	}
	return &barang, nil
}

func (r *BarangRepository) FindByKode(kode string) (*models.MasterBarang, error) {
	var barang models.MasterBarang
	err := r.db.Where("kode_barang = ?", kode).First(&barang).Error
	if err != nil {
		return nil, err
	}
	return &barang, nil
}

func (r *BarangRepository) GetAll(search string, page, limit int) ([]models.MasterBarang, int64, error) {
	var barangs []models.MasterBarang
	var total int64

	query := r.db.Model(&models.MasterBarang{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("nama_barang ILIKE ? OR kode_barang ILIKE ?", searchPattern, searchPattern)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("id ASC").Find(&barangs).Error

	return barangs, total, err
}

func (r *BarangRepository) GetAllWithStok(search string, page, limit int) ([]models.BarangWithStok, int64, error) {
	var results []models.BarangWithStok
	var total int64

	query := r.db.Table("master_barang").
		Select("master_barang.*, COALESCE(mstok.stok_akhir, 0) as stok_akhir").
		Joins("LEFT JOIN mstok ON mstok.barang_id = master_barang.id")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("master_barang.nama_barang ILIKE ? OR master_barang.kode_barang ILIKE ?", searchPattern, searchPattern)
	}

	r.db.Table("master_barang").Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("master_barang.id ASC").Scan(&results).Error

	return results, total, err
}
