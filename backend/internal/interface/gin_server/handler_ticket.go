package gin_server

import (
	"net/http"
	"strconv"

	"github.com/TicketX/backend/internal/interface/gin_server/middleware"
	"github.com/TicketX/backend/internal/use_case"
	"github.com/gin-gonic/gin"
)

// CreateTicketRequest represents the request body for creating a ticket
type CreateTicketRequest struct {
	Title       string  `json:"title" binding:"required"`
	Venue       string  `json:"venue" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Description string  `json:"description"`
}

// CreateTicket handles ticket creation
func (s *GinServer) CreateTicket(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := s.useCase.CreateTicket(c.Request.Context(), &use_case.CreateTicketParams{
		SellerID:    user.ID,
		Title:       req.Title,
		Venue:       req.Venue,
		Price:       req.Price,
		Category:    req.Category,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result.Ticket)
}

// GetTickets handles listing tickets
func (s *GinServer) GetTickets(c *gin.Context) {
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	title := c.Query("title")
	venue := c.Query("venue")
	category := c.Query("category")

	params := &use_case.ListTicketsParams{}
	if minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			params.MinPrice = &val
		}
	}
	if maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			params.MaxPrice = &val
		}
	}
	if title != "" {
		params.Title = &title
	}
	if venue != "" {
		params.Venue = &venue
	}
	if category != "" {
		params.Category = &category
	}

	result, err := s.useCase.ListTickets(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result.Tickets)
}

// GetTicketByID handles getting a ticket by ID
func (s *GinServer) GetTicketByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	result, err := s.useCase.GetTicketByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result.Ticket)
}

// GetMyTickets handles getting user's tickets
func (s *GinServer) GetMyTickets(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	result, err := s.useCase.GetMyTickets(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result.Tickets)
}

// DeleteTicket handles ticket deletion
func (s *GinServer) DeleteTicket(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	result, err := s.useCase.DeleteTicket(c.Request.Context(), id, user.ID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": result.Message})
}
