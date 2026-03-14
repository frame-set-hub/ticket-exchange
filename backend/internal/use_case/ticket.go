package use_case

import (
	"context"
	"errors"

	"github.com/TicketX/backend/internal/entity/ticket"
	"github.com/TicketX/backend/internal/repository/ticket_repository"
)

// CreateTicketParams contains parameters for creating a ticket
type CreateTicketParams struct {
	SellerID    uint
	Title       string
	Venue       string
	Price       float64
	Category    string
	Description string
}

// CreateTicketResult contains result of creating a ticket
type CreateTicketResult struct {
	Ticket *ticket.Ticket
}

// ListTicketsParams contains parameters for listing tickets
type ListTicketsParams struct {
	MinPrice *float64
	MaxPrice *float64
	Title    *string
	Venue    *string
	Category *string
}

// ListTicketsResult contains result of listing tickets
type ListTicketsResult struct {
	Tickets []*ticket.TicketWithSeller
}

// GetTicketByIDResult contains result of getting a ticket by ID
type GetTicketByIDResult struct {
	Ticket *ticket.TicketWithSeller
}

// GetMyTicketsResult contains result of getting user's tickets
type GetMyTicketsResult struct {
	Tickets []*ticket.Ticket
}

// DeleteTicketResult contains result of deleting a ticket
type DeleteTicketResult struct {
	Message string
}

// CreateTicket creates a new ticket
func (u *UseCase) CreateTicket(ctx context.Context, p *CreateTicketParams) (*CreateTicketResult, error) {
	ticketEntity := &ticket.Ticket{
		SellerID:    p.SellerID,
		Title:       p.Title,
		Venue:       p.Venue,
		Price:       p.Price,
		Category:    p.Category,
		Description: p.Description,
		Status:      ticket.TicketAvailable,
	}

	if err := u.ticketRepository.Create(ctx, ticketEntity); err != nil {
		return nil, errors.New("failed to create ticket")
	}

	return &CreateTicketResult{
		Ticket: ticketEntity,
	}, nil
}

// ListTickets lists tickets with optional filters
func (u *UseCase) ListTickets(ctx context.Context, p *ListTicketsParams) (*ListTicketsResult, error) {
	repoParams := &ticket_repository.ListTicketParams{}
	if p != nil {
		repoParams = &ticket_repository.ListTicketParams{
			MinPrice: p.MinPrice,
			MaxPrice: p.MaxPrice,
			Title:    p.Title,
			Venue:    p.Venue,
			Category: p.Category,
		}
	}

	tickets, err := u.ticketRepository.List(ctx, repoParams)
	if err != nil {
		return nil, errors.New("failed to fetch tickets")
	}

	return &ListTicketsResult{
		Tickets: tickets,
	}, nil
}

// GetTicketByID gets a ticket by ID
func (u *UseCase) GetTicketByID(ctx context.Context, id int) (*GetTicketByIDResult, error) {
	ticketEntity, err := u.ticketRepository.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	return &GetTicketByIDResult{
		Ticket: ticketEntity,
	}, nil
}

// GetMyTickets gets all tickets owned by a seller
func (u *UseCase) GetMyTickets(ctx context.Context, sellerID uint) (*GetMyTicketsResult, error) {
	tickets, err := u.ticketRepository.GetBySellerID(ctx, sellerID)
	if err != nil {
		return nil, errors.New("failed to fetch your tickets")
	}

	return &GetMyTicketsResult{
		Tickets: tickets,
	}, nil
}

// DeleteTicket deletes a ticket
func (u *UseCase) DeleteTicket(ctx context.Context, ticketID int, sellerID uint) (*DeleteTicketResult, error) {
	// Find ticket
	ticketEntity, err := u.ticketRepository.FindByID(ctx, ticketID)
	if err != nil {
		return nil, errors.New("ticket not found")
	}

	// Check ownership
	if ticketEntity.SellerID != sellerID {
		return nil, errors.New("you can only delete your own tickets")
	}

	// Check status
	if ticketEntity.Status != ticket.TicketAvailable {
		return nil, errors.New("cannot delete a ticket that is already in escrow")
	}

	// Delete ticket
	if err := u.ticketRepository.Delete(ctx, ticketID); err != nil {
		return nil, errors.New("failed to delete ticket")
	}

	return &DeleteTicketResult{
		Message: "Ticket deleted successfully",
	}, nil
}

// UpdateTicketStatus updates ticket status
func (u *UseCase) UpdateTicketStatus(ctx context.Context, ticketID int, status ticket.TicketStatus) error {
	return u.ticketRepository.UpdateStatus(ctx, ticketID, status)
}
