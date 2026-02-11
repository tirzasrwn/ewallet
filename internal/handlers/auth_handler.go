package handlers

import (
	"ewallet/internal/service"
	"ewallet/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"alice"`
	Email    string `json:"email" binding:"required,email" example:"alice@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"alice@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type LoginResponse struct {
	Token string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  interface{} `json:"user"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration Request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Registration failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", user.ToResponse())
}

// Login godoc
// @Summary Login user
// @Description Login with email and password to get JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err)
		return
	}

	response := LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}
