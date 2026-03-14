package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TicketX/backend/internal/entity/ticket"
	"github.com/TicketX/backend/internal/entity/transaction"
	"github.com/TicketX/backend/internal/entity/user"
	"github.com/TicketX/backend/internal/interface/gin_server"
	"github.com/TicketX/backend/internal/repository/message_repository"
	"github.com/TicketX/backend/internal/repository/ticket_repository"
	"github.com/TicketX/backend/internal/repository/transaction_repository"
	"github.com/TicketX/backend/internal/repository/user_repository"
	"github.com/TicketX/backend/internal/use_case"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database models for auto-migration (must match repository models)
type userModel struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string         `gorm:"uniqueIndex;not null"`
	Email     string         `gorm:"uniqueIndex;not null"`
	Password  string         `gorm:"not null"`
	Role      user.Role      `gorm:"type:varchar(20);default:'User'"`
}

type ticketModel struct {
	ID          uint               `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt     `gorm:"index"`
	SellerID    uint
	Title       string             `gorm:"not null"`
	Venue       string             `gorm:"not null"`
	Price       float64            `gorm:"not null"`
	Category    string             `gorm:"not null"`
	Description string
	Status      ticket.TicketStatus `gorm:"type:varchar(20);default:'Available'"`
}

type transactionModel struct {
	ID        uint                         `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	TicketID  uint
	BuyerID   uint
	SellerID  uint
	Status    transaction.TransactionStatus `gorm:"type:varchar(30);default:'Waiting_Ticket'"`
}

type messageModel struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	TransactionID uint
	SenderID      uint
	ReceiverID    uint
	Content       string
	AttachmentURL string
}

type config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	Port       string
	Debug      bool
	RequestLog bool
	AdminUserID uint
}

func initConfig() config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	return config{
		DBHost:     getEnv("DB_HOST"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME"),
		DBPort:     getEnv("DB_PORT"),
		Port:       getEnvOrDefault("PORT", "8080"),
		Debug:      getEnv("GIN_MODE") != "release",
		RequestLog: getEnv("GIN_MODE") != "release",
		AdminUserID: 1, // Default admin user ID
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func connectDB(cfg config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Database connection successfully opened")
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&userModel{},
		&ticketModel{},
		&transactionModel{},
		&messageModel{},
	)
}

func main() {
	cfg := initConfig()

	// Connect to database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate database
	if err := autoMigrate(db); err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	fmt.Println("Database auto-migrated successfully")

	// Initialize repositories
	userRepo := user_repository.NewPostgresDb(db)
	ticketRepo := ticket_repository.NewPostgresDb(db)
	txRepo := transaction_repository.NewPostgresDb(db)
	msgRepo := message_repository.NewPostgresDb(db)

	// Initialize use case
	useCase := use_case.New(
		use_case.Config{
			AdminUserID: cfg.AdminUserID,
		},
		userRepo,
		ticketRepo,
		txRepo,
		msgRepo,
	)

	// Initialize and start server
	server := gin_server.New(useCase, &gin_server.ServerConfig{
		Port:        cfg.Port,
		Debug:       cfg.Debug,
		RequestLog:  cfg.RequestLog,
	})

	server.SetupRoutes()

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
