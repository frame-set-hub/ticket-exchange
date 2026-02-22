package handlers

import (
	"net/http"
	"strconv"

	"github.com/TicketX/backend/internal/database"
	"github.com/TicketX/backend/internal/models"
	"github.com/gin-gonic/gin"
)

// Buyer initiates purchase
func CreateTransaction(c *gin.Context) {
	buyerID, _ := c.Get("user_id")

	ticketIDParam := c.Param("id")
	ticketID, err := strconv.Atoi(ticketIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var ticket models.Ticket
	if err := database.DB.First(&ticket, ticketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	if ticket.Status != models.TicketAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket is not available"})
		return
	}

	if ticket.SellerID == uint(buyerID.(float64)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot buy your own ticket"})
		return
	}

	// Update ticket status to prevent double booking
	ticket.Status = models.TicketPending
	database.DB.Save(&ticket)

	// Create Escrow Transaction
	tx := models.Transaction{
		TicketID: uint(ticketID),
		BuyerID:  uint(buyerID.(float64)),
		SellerID: ticket.SellerID,
		Status:   models.TxWaitingTicket,
	}

	if err := database.DB.Create(&tx).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate transaction"})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

// Seller provides the ticket
func SellerUploadTicket(c *gin.Context) {
	sellerID, _ := c.Get("user_id")
	txID := c.Param("id")

	var tx models.Transaction
	if err := database.DB.First(&tx, txID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if tx.SellerID != uint(sellerID.(float64)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	if tx.Status != models.TxWaitingTicket {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction is not in Waiting_Ticket state"})
		return
	}

	// Simulate upload
	tx.Status = models.TxWaitingPayment
	database.DB.Save(&tx)

	c.JSON(http.StatusOK, gin.H{"message": "Ticket uploaded to Escrow successfully", "transaction": tx})
}

// Buyer provides payment proof
func BuyerUploadPayment(c *gin.Context) {
	buyerID, _ := c.Get("user_id")
	txID := c.Param("id")

	var tx models.Transaction
	if err := database.DB.First(&tx, txID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if tx.BuyerID != uint(buyerID.(float64)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	if tx.Status != models.TxWaitingPayment {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction is not in Waiting_Payment state"})
		return
	}

	tx.Status = models.TxVerifying
	database.DB.Save(&tx)

	c.JSON(http.StatusOK, gin.H{"message": "Payment proof uploaded, awaiting Admin verification", "transaction": tx})
}

// Admin finalizes
func CompleteTransaction(c *gin.Context) {
	// Assumes Admin middleware protects this route
	txID := c.Param("id")

	var tx models.Transaction
	if err := database.DB.Preload("Ticket").First(&tx, txID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if tx.Status != models.TxVerifying {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction is not verifying"})
		return
	}

	// Move funds, transfer ticket theoretically happens here
	tx.Status = models.TxCompleted
	tx.Ticket.Status = models.TicketSold

	database.DB.Save(&tx)
	database.DB.Save(&tx.Ticket)

	c.JSON(http.StatusOK, gin.H{"message": "Transaction completed. Funds released and ticket transferred.", "transaction": tx})
}

func GetAllTransactions(c *gin.Context) {
	var txs []models.Transaction
	database.DB.Preload("Buyer").Preload("Seller").Preload("Ticket").Find(&txs)
	c.JSON(http.StatusOK, txs)
}
