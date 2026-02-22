package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TicketX/backend/internal/database"
	"github.com/TicketX/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTicketInput struct {
	Title       string  `json:"title" binding:"required"`
	Venue       string  `json:"venue" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Description string  `json:"description"`
}

func CreateTicket(c *gin.Context) {
	sellerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input CreateTicketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket := models.Ticket{
		SellerID:    uint(sellerID.(float64)), // JWT decodes numbers to float64
		Title:       input.Title,
		Venue:       input.Venue,
		Price:       input.Price,
		Category:    input.Category,
		Description: input.Description,
		Status:      models.TicketAvailable,
	}

	if err := database.DB.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func GetTickets(c *gin.Context) {
	// Query parameters for filtering
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	title := c.Query("title")
	venue := c.Query("venue")
	category := c.Query("category")

	query := database.DB.Model(&models.Ticket{}).Where("status = ?", models.TicketAvailable)

	if minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			query = query.Where("price >= ?", minPrice)
		}
	}
	if maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			query = query.Where("price <= ?", maxPrice)
		}
	}
	if title != "" {
		query = query.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	if venue != "" {
		query = query.Where("venue ILIKE ?", fmt.Sprintf("%%%s%%", venue))
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var tickets []models.Ticket
	if err := query.Preload("Seller", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username") // Preload seller with only safe fields
	}).Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func GetTicketByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	var ticket models.Ticket
	if err := database.DB.Preload("Seller").First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}
