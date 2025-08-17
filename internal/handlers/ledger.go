package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/obirikan/Ledger---Financial-Transaction-Management-API/internal/service"
	"github.com/shopspring/decimal"
)

type LedgerHandler struct {
	ledgerSvc service.LedgerService
}

func NewLedgerHandler(ledgerSvc service.LedgerService) *LedgerHandler {
	return &LedgerHandler{ledgerSvc: ledgerSvc}
}

func (h *LedgerHandler) ListAccounts(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id missing in context"})
		return
	}
	userID := userIDVal.(int)

	accounts, err := h.ledgerSvc.ListAccounts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

type transferRequest struct {
	FromAccountID int    `json:"from_account_id" binding:"required"`
	ToAccountID   int    `json:"to_account_id" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
}

func (h *LedgerHandler) TransferFunds(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id missing in context"})
		return
	}
	userID := userIDVal.(int)

	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}
	err = h.ledgerSvc.TransferFunds(c.Request.Context(), userID, req.FromAccountID, req.ToAccountID, amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transfer successful"})
}
