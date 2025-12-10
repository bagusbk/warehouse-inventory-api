package handlers

import (
	"net/http"
	"strconv"
	"warehouse/models"
	"warehouse/repositories"

	"github.com/gin-gonic/gin"
)

type BarangHandler struct {
	repo *repositories.BarangRepository
}

func NewBarangHandler() *BarangHandler {
	return &BarangHandler{repo: repositories.NewBarangRepository()}
}

func (h *BarangHandler) GetAll(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	barangs, total, err := h.repo.GetAll(search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to retrieve data", models.ErrInternalError))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMeta("Data retrieved successfully", barangs, &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}))
}

func (h *BarangHandler) GetAllWithStok(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	barangs, total, err := h.repo.GetAllWithStok(search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to retrieve data", models.ErrInternalError))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMeta("Data retrieved successfully", barangs, &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}))
}

func (h *BarangHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Invalid ID", models.ErrValidationError))
		return
	}

	barang, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponseMsg("Barang tidak ditemukan", models.ErrItemNotFound))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Data retrieved successfully", barang))
}

func (h *BarangHandler) Create(c *gin.Context) {
	var req models.BarangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseMsg(err.Error(), models.ErrValidationError))
		return
	}

	if _, err := h.repo.FindByKode(req.KodeBarang); err == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Kode barang sudah ada", models.ErrValidationError))
		return
	}

	barang := models.MasterBarang{
		KodeBarang: req.KodeBarang,
		NamaBarang: req.NamaBarang,
		Deskripsi:  req.Deskripsi,
		Satuan:     req.Satuan,
		HargaBeli:  req.HargaBeli,
		HargaJual:  req.HargaJual,
	}

	if err := h.repo.Create(&barang); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to create barang", models.ErrInternalError))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Barang created successfully", barang))
}

func (h *BarangHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Invalid ID", models.ErrValidationError))
		return
	}

	barang, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponseMsg("Barang tidak ditemukan", models.ErrItemNotFound))
		return
	}

	var req models.BarangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseMsg(err.Error(), models.ErrValidationError))
		return
	}

	if req.KodeBarang != barang.KodeBarang {
		if existing, _ := h.repo.FindByKode(req.KodeBarang); existing != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Kode barang sudah digunakan", models.ErrValidationError))
			return
		}
	}

	barang.KodeBarang = req.KodeBarang
	barang.NamaBarang = req.NamaBarang
	barang.Deskripsi = req.Deskripsi
	barang.Satuan = req.Satuan
	barang.HargaBeli = req.HargaBeli
	barang.HargaJual = req.HargaJual

	if err := h.repo.Update(barang); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to update barang", models.ErrInternalError))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Barang updated successfully", barang))
}

func (h *BarangHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Invalid ID", models.ErrValidationError))
		return
	}

	if _, err := h.repo.FindByID(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponseMsg("Barang tidak ditemukan", models.ErrItemNotFound))
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to delete barang", models.ErrInternalError))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Barang deleted successfully", nil))
}
