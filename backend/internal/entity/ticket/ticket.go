package ticket

import "time"

// Ticket represents the domain entity for ticket
type Ticket struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	SellerID    uint       `json:"seller_id"`
	Title       string     `json:"title"`
	Venue       string     `json:"venue"`
	Price       float64    `json:"price"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Status      TicketStatus `json:"status"`
}

// TicketWithSeller represents ticket with seller information
type TicketWithSeller struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	SellerID    uint       `json:"seller_id"`
	Seller      SellerInfo `json:"seller"`
	Title       string     `json:"title"`
	Venue       string     `json:"venue"`
	Price       float64    `json:"price"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Status      TicketStatus `json:"status"`
}

// SellerInfo represents seller information in ticket
type SellerInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// TicketStatus represents ticket statuses
type TicketStatus string

const (
	TicketAvailable TicketStatus = "Available"
	TicketPending   TicketStatus = "Pending"
	TicketSold      TicketStatus = "Sold"
)
