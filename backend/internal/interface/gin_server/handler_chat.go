package gin_server

import (
	"net/http"
	"strconv"

	"github.com/TicketX/backend/internal/interface/gin_server/middleware"
	"github.com/TicketX/backend/internal/use_case"
	"github.com/gin-gonic/gin"
)

// SendMessageRequest represents the request body for sending a message
type SendMessageRequest struct {
	Content       string `json:"content" binding:"required"`
	AttachmentURL string `json:"attachment_url"`
}

// GetMessages handles getting messages for a transaction
func (s *GinServer) GetMessages(c *gin.Context) {
	transactionIDStr := c.Param("transaction_id")
	transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	result, err := s.useCase.ListMessagesByTransactionID(c.Request.Context(), uint(transactionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result.Messages)
}

// SendMessage handles sending a message
func (s *GinServer) SendMessage(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	transactionIDStr := c.Param("transaction_id")
	transactionID, err := strconv.ParseUint(transactionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get transaction to find receiver
	txResult, err := s.useCase.GetTransactionByID(c.Request.Context(), uint(transactionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// Determine receiver (if sender is buyer, receiver is seller, and vice versa)
	receiverID := txResult.Transaction.BuyerID
	if user.ID == txResult.Transaction.BuyerID {
		receiverID = txResult.Transaction.SellerID
	}

	result, err := s.useCase.CreateMessage(c.Request.Context(), &use_case.CreateMessageParams{
		TransactionID: uint(transactionID),
		SenderID:      user.ID,
		ReceiverID:    receiverID,
		Content:       req.Content,
		AttachmentURL: req.AttachmentURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result.Message)
}
