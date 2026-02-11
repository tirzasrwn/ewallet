package handlers

import (
	"ewallet/internal/middleware"
	"ewallet/internal/service"
	"ewallet/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

type TransferRequest struct {
	ReceiverID uint    `json:"receiver_id" binding:"required,gt=0" example:"2"`
	Amount     float64 `json:"amount" binding:"required,gt=0" example:"50000"`
}

// Transfer godoc
// @Summary Transfer money to another user
// @Description Transfer funds from authenticated user's wallet to another user
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TransferRequest true "Transfer Request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/transactions/transfer [post]
func (h *TransactionHandler) Transfer(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	transaction, err := h.transactionService.Transfer(userID, req.ReceiverID, req.Amount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Transfer failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transfer successful", transaction.ToResponse())
}

// GetHistory godoc
// @Summary Get transaction history
// @Description Get authenticated user's transaction history
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit number of transactions" default(50)
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/transactions/history [get]
func (h *TransactionHandler) GetHistory(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Get limit from query parameter (default 50)
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	transactions, err := h.transactionService.GetHistory(userID, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve history", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction history retrieved successfully", transactions)
}
