package handlers

import (
	"net/http"
	"strconv"
	"warehouse/middleware"
	"warehouse/models"
	"warehouse/repositories"

	"github.com/gin-gonic/gin"
)

type PembelianHandler struct {
	repo *repositories.PembelianRepository
}

func NewPembelianHandler() *PembelianHandler {
	return &PembelianHandler{repo: repositories.NewPembelianRepository()}
}

func (h *PembelianHandler) Create(c *gin.Context) {
	var req models.PembelianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseMsg(err.Error(), models.ErrValidationError))
		return
	}

	if len(req.Details) == 0 {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseMsg("Detail pembelian tidak boleh kosong", models.ErrValidationError))
		return
	}

	userID := middleware.GetUserID(c)

	header, err := h.repo.Create(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg(err.Error(), models.ErrValidationError))
		return
	}

	result, _ := h.repo.GetByID(header.ID)

	response := h.formatResponse(result)

	c.JSON(http.StatusCreated, models.SuccessResponse("Pembelian berhasil dibuat", response))
}

func (h *PembelianHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	headers, total, err := h.repo.GetAll(page, limit, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to retrieve data", models.ErrInternalError))
		return
	}

	var response []models.PembelianResponse
	for _, header := range headers {
		response = append(response, h.formatResponse(&header))
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMeta("Data retrieved successfully", response, &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}))
}

func (h *PembelianHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Invalid ID", models.ErrValidationError))
		return
	}

	header, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponseMsg("Pembelian tidak ditemukan", models.ErrItemNotFound))
		return
	}

	response := h.formatResponse(header)

	c.JSON(http.StatusOK, models.SuccessResponse("Data retrieved successfully", response))
}

func (h *PembelianHandler) formatResponse(header *models.BeliHeader) models.PembelianResponse {
	var details []models.BeliDetailResponse
	for _, d := range header.Details {
		details = append(details, models.BeliDetailResponse{
			ID:       d.ID,
			BarangID: d.BarangID,
			Qty:      d.Qty,
			Harga:    d.Harga,
			Subtotal: d.Subtotal,
			Barang: models.BarangSimple{
				KodeBarang: d.Barang.KodeBarang,
				NamaBarang: d.Barang.NamaBarang,
				Satuan:     d.Barang.Satuan,
			},
		})
	}

	return models.PembelianResponse{
		Header: models.BeliHeaderResponse{
			ID:        header.ID,
			NoFaktur:  header.NoFaktur,
			Supplier:  header.Supplier,
			Total:     header.Total,
			UserID:    header.UserID,
			Status:    header.Status,
			CreatedAt: header.CreatedAt,
			User: models.UserResponse{
				ID:       header.User.ID,
				Username: header.User.Username,
				FullName: header.User.FullName,
			},
		},
		Details: details,
	}
}
