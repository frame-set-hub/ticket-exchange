package routes

import (
	"github.com/TicketX/backend/internal/handlers"
	"github.com/TicketX/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// CORS Middleware (simple for POC)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")

	// Auth
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Public Ticket Listings
	api.GET("/tickets", handlers.GetTickets)
	api.GET("/tickets/:id", handlers.GetTicketByID)

	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		// Tickets
		protected.POST("/tickets", handlers.CreateTicket)

		// Transactions
		protected.POST("/transactions/:id", handlers.CreateTransaction)
		protected.POST("/transactions/:id/upload-ticket", handlers.SellerUploadTicket)
		protected.POST("/transactions/:id/upload-payment", handlers.BuyerUploadPayment)

		// Chat WebSocket
		protected.GET("/chat", handlers.ConnectChat) // Need auth token inside URL query for WS ideally, or let it bypass middleware if sent during UPGRADE.

		// Admin only
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminRequired())
		{
			admin.GET("/transactions", handlers.GetAllTransactions)
			admin.POST("/transactions/:id/complete", handlers.CompleteTransaction)
		}
	}
}
