package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "User"
	RoleAdmin Role = "Admin"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      Role           `gorm:"type:varchar(20);default:'User'" json:"role"`
}

type TicketStatus string

const (
	TicketAvailable TicketStatus = "Available"
	TicketPending   TicketStatus = "Pending"
	TicketSold      TicketStatus = "Sold"
)

type Ticket struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	SellerID    uint           `json:"seller_id"`
	Seller      User           `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Title       string         `gorm:"not null" json:"title"`
	Venue       string         `gorm:"not null" json:"venue"`
	Price       float64        `gorm:"not null" json:"price"`
	Category    string         `gorm:"not null" json:"category"`
	Description string         `json:"description"`
	Status      TicketStatus   `gorm:"type:varchar(20);default:'Available'" json:"status"`
}

type TransactionStatus string

const (
	TxWaitingTicket  TransactionStatus = "Waiting_Ticket"
	TxWaitingPayment TransactionStatus = "Waiting_Payment"
	TxVerifying      TransactionStatus = "Verifying"
	TxCompleted      TransactionStatus = "Completed"
	TxCancelled      TransactionStatus = "Cancelled"
)

type Transaction struct {
	ID        uint              `gorm:"primarykey" json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	TicketID  uint              `json:"ticket_id"`
	Ticket    Ticket            `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	BuyerID   uint              `json:"buyer_id"`
	Buyer     User              `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	SellerID  uint              `json:"seller_id"`
	Seller    User              `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Status    TransactionStatus `gorm:"type:varchar(30);default:'Waiting_Ticket'" json:"status"`
}

type Message struct {
	ID            uint        `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time   `json:"created_at"`
	TransactionID uint        `json:"transaction_id"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID" json:"-"`
	SenderID      uint        `json:"sender_id"`
	ReceiverID    uint        `json:"receiver_id"` // Usually Admin ID
	Content       string      `json:"content"`
	AttachmentURL string      `json:"attachment_url"`
}
