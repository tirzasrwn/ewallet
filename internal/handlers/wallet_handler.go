package handlers

import (
	"ewallet/internal/middleware"
	"ewallet/internal/service"
	"ewallet/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService service.WalletService
}

func NewWalletHandler(walletService service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0" example:"100000"`
}

// GetBalance godoc
// @Summary Get wallet balance
// @Description Get authenticated user's wallet balance
// @Tags Wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/wallets/balance [get]
func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	wallet, err := h.walletService.GetBalance(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Wallet not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Balance retrieved successfully", wallet)
}

// TopUp godoc
// @Summary Top up wallet balance
// @Description Add funds to authenticated user's wallet
// @Tags Wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TopUpRequest true "Top Up Request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/wallets/topup [post]
func (h *WalletHandler) TopUp(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	wallet, err := h.walletService.TopUp(userID, req.Amount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Top up failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Top up successful", wallet)
}
