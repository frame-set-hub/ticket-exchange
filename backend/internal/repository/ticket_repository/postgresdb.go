package ticket_repository

import (
	"context"
	"errors"
	"time"

	"github.com/TicketX/backend/internal/entity/ticket"
	"gorm.io/gorm"
)

// pgDb implements TicketRepository interface with PostgreSQL
type pgDb struct {
	db *gorm.DB
}

// ticketModel represents the database model for ticket
type ticketModel struct {
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	SellerID    uint
	Title       string `gorm:"not null"`
	Venue       string `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Category    string `gorm:"not null"`
	Description string
	Status      ticket.TicketStatus `gorm:"type:varchar(20);default:'Available'"`
}

// sellerModel represents seller information
type sellerModel struct {
	ID       uint `gorm:"primarykey"`
	Username string
}

// NewPostgresDb creates a new PostgreSQL ticket repository
func NewPostgresDb(db *gorm.DB) TicketRepository {
	return &pgDb{db: db}
}

// Create creates a new ticket
func (r *pgDb) Create(ctx context.Context, t *ticket.Ticket) error {
	ticketModel := &ticketModel{
		SellerID:    t.SellerID,
		Title:       t.Title,
		Venue:       t.Venue,
		Price:       t.Price,
		Category:    t.Category,
		Description: t.Description,
		Status:      t.Status,
	}

	if err := r.db.WithContext(ctx).Create(ticketModel).Error; err != nil {
		return err
	}

	t.ID = ticketModel.ID
	return nil
}

// List lists tickets with optional filters
func (r *pgDb) List(ctx context.Context, params *ListTicketParams) ([]*ticket.TicketWithSeller, error) {
	query := r.db.WithContext(ctx).Model(&ticketModel{}).Where("status = ?", ticket.TicketAvailable)

	if params != nil {
		if params.MinPrice != nil {
			query = query.Where("price >= ?", *params.MinPrice)
		}
		if params.MaxPrice != nil {
			query = query.Where("price <= ?", *params.MaxPrice)
		}
		if params.Title != nil {
			query = query.Where("title ILIKE ?", "%"+*params.Title+"%")
		}
		if params.Venue != nil {
			query = query.Where("venue ILIKE ?", "%"+*params.Venue+"%")
		}
		if params.Category != nil {
			query = query.Where("category = ?", *params.Category)
		}
	}

	var tickets []ticketModel

	if err := query.Find(&tickets).Error; err != nil {
		return nil, err
	}

	// Transform to domain entities
	result := make([]*ticket.TicketWithSeller, len(tickets))
	for i, t := range tickets {
		result[i] = &ticket.TicketWithSeller{
			ID:          t.ID,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
			SellerID:    t.SellerID,
			Title:       t.Title,
			Venue:       t.Venue,
			Price:       t.Price,
			Category:    t.Category,
			Description: t.Description,
			Status:      t.Status,
		}
	}

	return result, nil
}

// GetByID gets a ticket by ID
func (r *pgDb) GetByID(ctx context.Context, id int) (*ticket.TicketWithSeller, error) {
	var ticketModel ticketModel

	err := r.db.WithContext(ctx).First(&ticketModel, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}

	return &ticket.TicketWithSeller{
		ID:          ticketModel.ID,
		CreatedAt:   ticketModel.CreatedAt,
		UpdatedAt:   ticketModel.UpdatedAt,
		SellerID:    ticketModel.SellerID,
		Title:       ticketModel.Title,
		Venue:       ticketModel.Venue,
		Price:       ticketModel.Price,
		Category:    ticketModel.Category,
		Description: ticketModel.Description,
		Status:      ticketModel.Status,
	}, nil
}

// GetBySellerID gets all tickets by seller ID
func (r *pgDb) GetBySellerID(ctx context.Context, sellerID uint) ([]*ticket.Ticket, error) {
	var tickets []ticketModel

	if err := r.db.WithContext(ctx).Where("seller_id = ?", sellerID).Find(&tickets).Error; err != nil {
		return nil, err
	}

	result := make([]*ticket.Ticket, len(tickets))
	for i, t := range tickets {
		result[i] = &ticket.Ticket{
			ID:          t.ID,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
			SellerID:    t.SellerID,
			Title:       t.Title,
			Venue:       t.Venue,
			Price:       t.Price,
			Category:    t.Category,
			Description: t.Description,
			Status:      t.Status,
		}
	}

	return result, nil
}

// Delete deletes a ticket
func (r *pgDb) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&ticketModel{}, id).Error
}

// FindByID finds a ticket by ID (without seller)
func (r *pgDb) FindByID(ctx context.Context, id int) (*ticket.Ticket, error) {
	var ticketModel ticketModel

	if err := r.db.WithContext(ctx).First(&ticketModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}

	return &ticket.Ticket{
		ID:          ticketModel.ID,
		CreatedAt:   ticketModel.CreatedAt,
		UpdatedAt:   ticketModel.UpdatedAt,
		SellerID:    ticketModel.SellerID,
		Title:       ticketModel.Title,
		Venue:       ticketModel.Venue,
		Price:       ticketModel.Price,
		Category:    ticketModel.Category,
		Description: ticketModel.Description,
		Status:      ticketModel.Status,
	}, nil
}

// UpdateStatus updates ticket status
func (r *pgDb) UpdateStatus(ctx context.Context, id int, status ticket.TicketStatus) error {
	return r.db.WithContext(ctx).Model(&ticketModel{}).Where("id = ?", id).Update("status", status).Error
}
