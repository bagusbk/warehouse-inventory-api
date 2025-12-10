package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"warehouse/middleware"
	"warehouse/models"
	"warehouse/repositories"

	"github.com/gin-gonic/gin"
)

type PenjualanHandler struct {
	repo *repositories.PenjualanRepository
}

func NewPenjualanHandler() *PenjualanHandler {
	return &PenjualanHandler{repo: repositories.NewPenjualanRepository()}
}

func (h *PenjualanHandler) Create(c *gin.Context) {
	var req models.PenjualanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseMsg(err.Error(), models.ErrValidationError))
		return
	}

	if len(req.Details) == 0 {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseMsg("Detail penjualan tidak boleh kosong", models.ErrValidationError))
		return
	}

	userID := middleware.GetUserID(c)

	header, err := h.repo.Create(&req, userID)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "stok tidak mencukupi") {
			c.JSON(http.StatusBadRequest, models.ErrorResponseMsg(errMsg, models.ErrInsufficientStock))
			return
		}
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg(errMsg, models.ErrValidationError))
		return
	}

	result, _ := h.repo.GetByID(header.ID)

	response := h.formatResponse(result)

	c.JSON(http.StatusCreated, models.SuccessResponse("Penjualan berhasil dibuat", response))
}

func (h *PenjualanHandler) GetAll(c *gin.Context) {
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

	var response []models.PenjualanResponse
	for _, header := range headers {
		response = append(response, h.formatResponse(&header))
	}

	c.JSON(http.StatusOK, models.SuccessResponseWithMeta("Data retrieved successfully", response, &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}))
}

func (h *PenjualanHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Invalid ID", models.ErrValidationError))
		return
	}

	header, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponseMsg("Penjualan tidak ditemukan", models.ErrItemNotFound))
		return
	}

	response := h.formatResponse(header)

	c.JSON(http.StatusOK, models.SuccessResponse("Data retrieved successfully", response))
}

func (h *PenjualanHandler) formatResponse(header *models.JualHeader) models.PenjualanResponse {
	var details []models.JualDetailResponse
	for _, d := range header.Details {
		details = append(details, models.JualDetailResponse{
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

	return models.PenjualanResponse{
		Header: models.JualHeaderResponse{
			ID:        header.ID,
			NoFaktur:  header.NoFaktur,
			Customer:  header.Customer,
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
