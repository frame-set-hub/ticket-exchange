package gin_server

import (
	"net/http"
	"strconv"

	"github.com/TicketX/backend/internal/entity/transaction"
	"github.com/TicketX/backend/internal/interface/gin_server/middleware"
	"github.com/TicketX/backend/internal/use_case"
	"github.com/gin-gonic/gin"
)

// CreateTransactionRequest represents the request body for creating a transaction
type CreateTransactionRequest struct {
	TicketID uint `json:"ticket_id" binding:"required"`
}

// UpdateTransactionStatusRequest represents the request body for updating transaction status
type UpdateTransactionStatusRequest struct {
	Status transaction.TransactionStatus `json:"status" binding:"required"`
}

// CreateTransaction handles transaction creation
func (s *GinServer) CreateTransaction(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get ticket to find seller
	ticketResult, err := s.useCase.GetTicketByID(c.Request.Context(), int(req.TicketID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	result, err := s.useCase.CreateTransaction(c.Request.Context(), &use_case.CreateTransactionParams{
		TicketID: req.TicketID,
		BuyerID:  user.ID,
		SellerID: ticketResult.Ticket.SellerID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result.Transaction)
}

// ListTransactions handles listing all transactions (admin only)
func (s *GinServer) ListTransactions(c *gin.Context) {
	result, err := s.useCase.ListTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result.Transactions)
}

// ListMyTransactions handles listing user's transactions
func (s *GinServer) ListMyTransactions(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	result, err := s.useCase.ListTransactionsByBuyerID(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result.Transactions)
}

// GetTransactionByID handles getting a transaction by ID
func (s *GinServer) GetTransactionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	result, err := s.useCase.GetTransactionByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result.Transaction)
}

// UpdateTransactionStatus handles updating transaction status
func (s *GinServer) UpdateTransactionStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var req UpdateTransactionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.useCase.UpdateTransactionStatus(c.Request.Context(), &use_case.UpdateTransactionStatusParams{
		TransactionID: uint(id),
		Status:        req.Status,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction status updated successfully"})
}

// AdminListTransactions handles listing all transactions (admin only)
func (s *GinServer) AdminListTransactions(c *gin.Context) {
	result, err := s.useCase.ListTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result.Transactions)
}

// AdminUpdateTransactionStatus handles updating transaction status (admin only)
func (s *GinServer) AdminUpdateTransactionStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var req UpdateTransactionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.useCase.UpdateTransactionStatus(c.Request.Context(), &use_case.UpdateTransactionStatusParams{
		TransactionID: uint(id),
		Status:        req.Status,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction status updated successfully"})
}
