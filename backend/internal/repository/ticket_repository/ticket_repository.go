package ticket_repository

import (
	"context"

	"github.com/TicketX/backend/internal/entity/ticket"
)

// ListTicketParams defines parameters for listing tickets
type ListTicketParams struct {
	MinPrice *float64
	MaxPrice *float64
	Title    *string
	Venue    *string
	Category *string
}

// TicketRepository defines the interface for ticket data operations
type TicketRepository interface {
	Create(ctx context.Context, t *ticket.Ticket) error
	List(ctx context.Context, params *ListTicketParams) ([]*ticket.TicketWithSeller, error)
	GetByID(ctx context.Context, id int) (*ticket.TicketWithSeller, error)
	GetBySellerID(ctx context.Context, sellerID uint) ([]*ticket.Ticket, error)
	Delete(ctx context.Context, id int) error
	FindByID(ctx context.Context, id int) (*ticket.Ticket, error)
	UpdateStatus(ctx context.Context, id int, status ticket.TicketStatus) error
}
