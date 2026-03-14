package transaction

import (
	"time"

	"github.com/TicketX/backend/internal/entity/ticket"
	"github.com/TicketX/backend/internal/entity/user"
)

// Transaction represents the domain entity for transaction
type Transaction struct {
	ID        uint            `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	TicketID  uint            `json:"ticket_id"`
	BuyerID   uint            `json:"buyer_id"`
	SellerID  uint            `json:"seller_id"`
	Status    TransactionStatus `json:"status"`
}

// TransactionWithDetails represents transaction with full details
type TransactionWithDetails struct {
	ID        uint              `json:"id"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	TicketID  uint              `json:"ticket_id"`
	Ticket    *ticket.Ticket    `json:"ticket,omitempty"`
	BuyerID   uint              `json:"buyer_id"`
	Buyer     *user.User        `json:"buyer,omitempty"`
	SellerID  uint              `json:"seller_id"`
	Seller    *user.User        `json:"seller,omitempty"`
	Status    TransactionStatus `json:"status"`
}

// TransactionStatus represents transaction statuses
type TransactionStatus string

const (
	TxWaitingTicket  TransactionStatus = "Waiting_Ticket"
	TxWaitingPayment TransactionStatus = "Waiting_Payment"
	TxVerifying      TransactionStatus = "Verifying"
	TxCompleted      TransactionStatus = "Completed"
	TxCancelled      TransactionStatus = "Cancelled"
)
