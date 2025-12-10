package handlers

import (
	"net/http"
	"warehouse/models"
	"warehouse/repositories"

	"github.com/gin-gonic/gin"
)

type LaporanHandler struct {
	stokRepo      *repositories.StokRepository
	pembelianRepo *repositories.PembelianRepository
	penjualanRepo *repositories.PenjualanRepository
}

func NewLaporanHandler() *LaporanHandler {
	return &LaporanHandler{
		stokRepo:      repositories.NewStokRepository(),
		pembelianRepo: repositories.NewPembelianRepository(),
		penjualanRepo: repositories.NewPenjualanRepository(),
	}
}

func (h *LaporanHandler) LaporanStok(c *gin.Context) {
	stoks, err := h.stokRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to retrieve data", models.ErrInternalError))
		return
	}

	var response []models.StokResponse
	var totalNilaiStok float64

	for _, stok := range stoks {
		nilaiStok := float64(stok.StokAkhir) * stok.Barang.HargaJual
		totalNilaiStok += nilaiStok

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

	c.JSON(http.StatusOK, models.SuccessResponse("Laporan stok retrieved", gin.H{
		"data":             response,
		"total_nilai_stok": totalNilaiStok,
	}))
}

func (h *LaporanHandler) LaporanPembelian(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	headers, total, err := h.pembelianRepo.GetAll(1, 1000, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to retrieve data", models.ErrInternalError))
		return
	}

	var totalPembelian float64
	for _, h := range headers {
		totalPembelian += h.Total
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Laporan pembelian retrieved", gin.H{
		"data":            headers,
		"total_transaksi": total,
		"total_pembelian": totalPembelian,
	}))
}

func (h *LaporanHandler) LaporanPenjualan(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	headers, total, err := h.penjualanRepo.GetAll(1, 1000, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to retrieve data", models.ErrInternalError))
		return
	}

	var totalPenjualan float64
	for _, h := range headers {
		totalPenjualan += h.Total
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Laporan penjualan retrieved", gin.H{
		"data":            headers,
		"total_transaksi": total,
		"total_penjualan": totalPenjualan,
	}))
}
