package transaction_repository

import (
	"context"
	"errors"
	"time"

	"github.com/TicketX/backend/internal/entity/transaction"
	"gorm.io/gorm"
)

// pgDb implements TransactionRepository interface with PostgreSQL
type pgDb struct {
	db *gorm.DB
}

// transactionModel represents the database model for transaction
type transactionModel struct {
	ID        uint                  `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	TicketID  uint
	BuyerID   uint
	SellerID  uint
	Status    transaction.TransactionStatus `gorm:"type:varchar(30);default:'Waiting_Ticket'"`
}

// NewPostgresDb creates a new PostgreSQL transaction repository
func NewPostgresDb(db *gorm.DB) TransactionRepository {
	return &pgDb{db: db}
}

// Create creates a new transaction
func (r *pgDb) Create(ctx context.Context, t *transaction.Transaction) error {
	txModel := &transactionModel{
		TicketID: t.TicketID,
		BuyerID:  t.BuyerID,
		SellerID: t.SellerID,
		Status:   t.Status,
	}

	if err := r.db.WithContext(ctx).Create(txModel).Error; err != nil {
		return err
	}

	t.ID = txModel.ID
	return nil
}

// List lists all transactions
func (r *pgDb) List(ctx context.Context) ([]*transaction.TransactionWithDetails, error) {
	var txModels []transactionModel

	if err := r.db.WithContext(ctx).Find(&txModels).Error; err != nil {
		return nil, err
	}

	result := make([]*transaction.TransactionWithDetails, len(txModels))
	for i, tx := range txModels {
		result[i] = &transaction.TransactionWithDetails{
			ID:        tx.ID,
			CreatedAt: tx.CreatedAt,
			UpdatedAt: tx.UpdatedAt,
			TicketID:  tx.TicketID,
			BuyerID:   tx.BuyerID,
			SellerID:  tx.SellerID,
			Status:    tx.Status,
		}
	}

	return result, nil
}

// GetByID gets a transaction by ID
func (r *pgDb) GetByID(ctx context.Context, id uint) (*transaction.TransactionWithDetails, error) {
	var txModel transactionModel

	err := r.db.WithContext(ctx).First(&txModel, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	return &transaction.TransactionWithDetails{
		ID:        txModel.ID,
		CreatedAt: txModel.CreatedAt,
		UpdatedAt: txModel.UpdatedAt,
		TicketID:  txModel.TicketID,
		BuyerID:   txModel.BuyerID,
		SellerID:  txModel.SellerID,
		Status:    txModel.Status,
	}, nil
}

// GetByTicketID gets a transaction by ticket ID
func (r *pgDb) GetByTicketID(ctx context.Context, ticketID uint) (*transaction.TransactionWithDetails, error) {
	var txModel transactionModel

	err := r.db.WithContext(ctx).Where("ticket_id = ?", ticketID).First(&txModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	return &transaction.TransactionWithDetails{
		ID:        txModel.ID,
		CreatedAt: txModel.CreatedAt,
		UpdatedAt: txModel.UpdatedAt,
		TicketID:  txModel.TicketID,
		BuyerID:   txModel.BuyerID,
		SellerID:  txModel.SellerID,
		Status:    txModel.Status,
	}, nil
}

// UpdateStatus updates transaction status
func (r *pgDb) UpdateStatus(ctx context.Context, id uint, status transaction.TransactionStatus) error {
	return r.db.WithContext(ctx).Model(&transactionModel{}).Where("id = ?", id).Update("status", status).Error
}

// ListByBuyerID lists transactions by buyer ID
func (r *pgDb) ListByBuyerID(ctx context.Context, buyerID uint) ([]*transaction.TransactionWithDetails, error) {
	var txModels []transactionModel

	if err := r.db.WithContext(ctx).
		Where("buyer_id = ?", buyerID).
		Find(&txModels).Error; err != nil {
		return nil, err
	}

	result := make([]*transaction.TransactionWithDetails, len(txModels))
	for i, tx := range txModels {
		result[i] = &transaction.TransactionWithDetails{
			ID:        tx.ID,
			CreatedAt: tx.CreatedAt,
			UpdatedAt: tx.UpdatedAt,
			TicketID:  tx.TicketID,
			BuyerID:   tx.BuyerID,
			SellerID:  tx.SellerID,
			Status:    tx.Status,
		}
	}

	return result, nil
}

// ListBySellerID lists transactions by seller ID
func (r *pgDb) ListBySellerID(ctx context.Context, sellerID uint) ([]*transaction.TransactionWithDetails, error) {
	var txModels []transactionModel

	if err := r.db.WithContext(ctx).
		Where("seller_id = ?", sellerID).
		Find(&txModels).Error; err != nil {
		return nil, err
	}

	result := make([]*transaction.TransactionWithDetails, len(txModels))
	for i, tx := range txModels {
		result[i] = &transaction.TransactionWithDetails{
			ID:        tx.ID,
			CreatedAt: tx.CreatedAt,
			UpdatedAt: tx.UpdatedAt,
			TicketID:  tx.TicketID,
			BuyerID:   tx.BuyerID,
			SellerID:  tx.SellerID,
			Status:    tx.Status,
		}
	}

	return result, nil
}
