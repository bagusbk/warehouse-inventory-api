package handlers

import (
	"net/http"
	"strconv"
	"warehouse/models"
	"warehouse/repositories"

	"github.com/gin-gonic/gin"
)

type StokHandler struct {
	repo *repositories.StokRepository
}

func NewStokHandler() *StokHandler {
	return &StokHandler{repo: repositories.NewStokRepository()}
}

func (h *StokHandler) GetAll(c *gin.Context) {
	stoks, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to retrieve data", models.ErrInternalError))
		return
	}

	var response []models.StokResponse
	for _, stok := range stoks {
		response = append(response, models.StokResponse{
			ID:        stok.ID,
			BarangID:  stok.BarangID,
			StokAkhir: stok.StokAkhir,
			UpdatedAt: stok.UpdatedAt,
			Barang: models.BarangSimple{
				KodeBarang: stok.Barang.KodeBarang,
				NamaBarang: stok.Barang.NamaBarang,
				Satuan:     stok.Barang.Satuan,
				HargaJual:  stok.Barang.HargaJual,
			},
		})
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Data retrieved successfully", response))
}

func (h *StokHandler) GetByBarang(c *gin.Context) {
	barangID, err := strconv.ParseUint(c.Param("barang_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Invalid barang ID", models.ErrValidationError))
		return
	}

	stok, err := h.repo.GetByBarangID(uint(barangID))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponseMsg("Stok tidak ditemukan untuk barang ini", models.ErrItemNotFound))
		return
	}

	response := models.StokResponse{
		ID:        stok.ID,
		BarangID:  stok.BarangID,
		StokAkhir: stok.StokAkhir,
		UpdatedAt: stok.UpdatedAt,
		Barang: models.BarangSimple{
			KodeBarang: stok.Barang.KodeBarang,
			NamaBarang: stok.Barang.NamaBarang,
			Satuan:     stok.Barang.Satuan,
			HargaJual:  stok.Barang.HargaJual,
		},
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Data retrieved successfully", response))
}

func (h *StokHandler) GetHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	history, total, err := h.repo.GetAllHistory(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to retrieve data", models.ErrInternalError))
		return
	}

	var response []models.HistoryStokResponse
	for _, h := range history {
		response = append(response, models.HistoryStokResponse{
			ID:             h.ID,
			BarangID:       h.BarangID,
			UserID:         h.UserID,
			JenisTransaksi: h.JenisTransaksi,
			Jumlah:         h.Jumlah,
			StokSebelum:    h.StokSebelum,
			StokSesudah:    h.StokSesudah,
			Keterangan:     h.Keterangan,
			CreatedAt:      h.CreatedAt,
			Barang: models.BarangSimple{
				KodeBarang: h.Barang.KodeBarang,
				NamaBarang: h.Barang.NamaBarang,
			},
			User: models.UserResponse{
				ID:       h.User.ID,
				Username: h.User.Username,
				FullName: h.User.FullName,
			},
		})
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMeta("Data retrieved successfully", response, &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}))
}

func (h *StokHandler) GetHistoryByBarang(c *gin.Context) {
	barangID, err := strconv.ParseUint(c.Param("barang_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Invalid barang ID", models.ErrValidationError))
		return
	}

	history, err := h.repo.GetHistoryByBarang(uint(barangID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to retrieve data", models.ErrInternalError))
		return
	}

	var response []models.HistoryStokResponse
	for _, h := range history {
		response = append(response, models.HistoryStokResponse{
			ID:             h.ID,
			BarangID:       h.BarangID,
			UserID:         h.UserID,
			JenisTransaksi: h.JenisTransaksi,
			Jumlah:         h.Jumlah,
			StokSebelum:    h.StokSebelum,
			StokSesudah:    h.StokSesudah,
			Keterangan:     h.Keterangan,
			CreatedAt:      h.CreatedAt,
			Barang: models.BarangSimple{
				KodeBarang: h.Barang.KodeBarang,
				NamaBarang: h.Barang.NamaBarang,
			},
			User: models.UserResponse{
				ID:       h.User.ID,
				Username: h.User.Username,
				FullName: h.User.FullName,
			},
		})
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Data retrieved successfully", response))
}
