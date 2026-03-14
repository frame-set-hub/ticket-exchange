package gin_server

import (
	"fmt"
	"net/http"

	"github.com/TicketX/backend/internal/interface/gin_server/middleware"
	"github.com/TicketX/backend/internal/use_case"
	"github.com/gin-gonic/gin"
)

// ServerConfig contains configuration for the Gin server
type ServerConfig struct {
	Port        string
	Debug       bool
	RequestLog  bool
}

// GinServer represents the HTTP server
type GinServer struct {
	useCase  *use_case.UseCase
	config   *ServerConfig
	router   *gin.Engine
	hub      *Hub
}

// New creates a new GinServer instance
func New(useCase *use_case.UseCase, config *ServerConfig) *GinServer {
	gin.SetMode(gin.ReleaseMode)
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// CORS — must be before any other middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	hub := NewHub()
	go hub.Run()

	return &GinServer{
		useCase:  useCase,
		config:   config,
		router:   router,
		hub:      hub,
	}
}

// SetupRoutes sets up all routes
func (s *GinServer) SetupRoutes() {
	// Public routes
	auth := s.router.Group("/api/auth")
	{
		auth.POST("/register", s.Register)
		auth.POST("/login", s.Login)
	}

	// Protected routes
	protected := s.router.Group("/api")
	protected.Use(middleware.Authenticate())
	{
		// Ticket routes
		tickets := protected.Group("/tickets")
		{
			tickets.POST("", s.CreateTicket)
			tickets.GET("", s.GetTickets)
			tickets.GET("/my", s.GetMyTickets)
			tickets.GET("/:id", s.GetTicketByID)
			tickets.DELETE("/:id", s.DeleteTicket)
		}

		// Transaction routes
		transactions := protected.Group("/transactions")
		{
			transactions.POST("", s.CreateTransaction)
			transactions.GET("", s.ListTransactions)
			transactions.GET("/my", s.ListMyTransactions)
			transactions.GET("/by-ticket/:ticket_id", s.GetTransactionByTicketID)
			transactions.GET("/:id", s.GetTransactionByID)
			transactions.POST("/:id/status", s.UpdateTransactionStatus)
		}

		// Chat routes
		chat := protected.Group("/chat")
		{
			chat.GET("/transactions/:transaction_id/messages", s.GetMessages)
			chat.POST("/transactions/:transaction_id/messages", s.SendMessage)
		}

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("/transactions", s.AdminListTransactions)
			admin.POST("/transactions/:id/status", s.AdminUpdateTransactionStatus)
		}
	}

	// WebSocket chat (auth via query param)
	s.router.GET("/api/chat/ws/:transaction_id", s.HandleWebSocket)

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// Start starts the HTTP server
func (s *GinServer) Start() error {
	addr := fmt.Sprintf(":%s", s.config.Port)
	fmt.Printf("Server starting on port %s...\n", s.config.Port)
	return s.router.Run(addr)
}
