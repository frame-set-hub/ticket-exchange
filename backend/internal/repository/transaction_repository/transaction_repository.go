package transaction_repository

import (
	"context"

	"github.com/TicketX/backend/internal/entity/transaction"
)

// TransactionRepository defines the interface for transaction data operations
type TransactionRepository interface {
	Create(ctx context.Context, t *transaction.Transaction) error
	List(ctx context.Context) ([]*transaction.TransactionWithDetails, error)
	GetByID(ctx context.Context, id uint) (*transaction.TransactionWithDetails, error)
	UpdateStatus(ctx context.Context, id uint, status transaction.TransactionStatus) error
	ListByBuyerID(ctx context.Context, buyerID uint) ([]*transaction.TransactionWithDetails, error)
	ListBySellerID(ctx context.Context, sellerID uint) ([]*transaction.TransactionWithDetails, error)
}
