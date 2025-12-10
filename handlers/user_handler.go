package handlers

import (
	"net/http"
	"strings"
	"warehouse/middleware"
	"warehouse/models"
	"warehouse/repositories"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo *repositories.UserRepository
}

func NewUserHandler() *UserHandler {
	return &UserHandler{repo: repositories.NewUserRepository()}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseMsg(err.Error(), models.ErrValidationError))
		return
	}

	if _, err := h.repo.FindByUsername(req.Username); err == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Username already exists", models.ErrValidationError))
		return
	}

	if _, err := h.repo.FindByEmail(req.Email); err == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponseMsg("Email already exists", models.ErrValidationError))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to hash password", models.ErrInternalError))
		return
	}

	role := "staff"
	if req.Role == "admin" {
		role = "admin"
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		FullName: req.FullName,
		Role:     role,
	}

	if err := h.repo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to create user", models.ErrInternalError))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("User registered successfully", models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Role:     user.Role,
	}))
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, models.ErrorResponseMsg(err.Error(), models.ErrValidationError))
		return
	}
	user, err := h.repo.FindByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponseMsg("Invalid username", models.ErrUnauthorized))
		return
	}

	isBcrypt := (strings.HasPrefix(user.Password, "$2a$") ||
		strings.HasPrefix(user.Password, "$2b$") ||
		strings.HasPrefix(user.Password, "$2y$")) &&
		len(user.Password) == 60

	if isBcrypt {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseMsg("Invalid username or password", models.ErrUnauthorized))
			return
		}
	} else {
		if user.Password != req.Password {
			c.JSON(http.StatusUnauthorized, models.ErrorResponseMsg("Invalid username or password", models.ErrUnauthorized))
			return
		}
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponseMsg("Failed to generate token", models.ErrInternalError))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Login successful", models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}))
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.repo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponseMsg("User not found", models.ErrItemNotFound))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Profile retrieved", models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Role:     user.Role,
	}))
}
