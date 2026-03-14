package use_case

import (
	"context"
	"errors"

	"github.com/TicketX/backend/internal/entity/ticket"
	"github.com/TicketX/backend/internal/entity/transaction"
)

// CreateTransactionParams contains parameters for creating a transaction
type CreateTransactionParams struct {
	TicketID uint
	BuyerID  uint
	SellerID uint
}

// CreateTransactionResult contains result of creating a transaction
type CreateTransactionResult struct {
	Transaction *transaction.Transaction
}

// ListTransactionsResult contains result of listing transactions
type ListTransactionsResult struct {
	Transactions []*transaction.TransactionWithDetails
}

// GetTransactionByIDResult contains result of getting a transaction by ID
type GetTransactionByIDResult struct {
	Transaction *transaction.TransactionWithDetails
}

// UpdateTransactionStatusParams contains parameters for updating transaction status
type UpdateTransactionStatusParams struct {
	TransactionID uint
	Status        transaction.TransactionStatus
}

// CreateTransaction creates a new transaction
func (u *UseCase) CreateTransaction(ctx context.Context, p *CreateTransactionParams) (*CreateTransactionResult, error) {
	txEntity := &transaction.Transaction{
		TicketID: p.TicketID,
		BuyerID:  p.BuyerID,
		SellerID: p.SellerID,
		Status:   transaction.TxWaitingTicket,
	}

	if err := u.transactionRepository.Create(ctx, txEntity); err != nil {
		return nil, errors.New("failed to create transaction")
	}

	// Update ticket status to pending
	if err := u.UpdateTicketStatus(ctx, int(p.TicketID), ticket.TicketPending); err != nil {
		return nil, errors.New("failed to update ticket status")
	}

	return &CreateTransactionResult{
		Transaction: txEntity,
	}, nil
}

// ListTransactions lists all transactions (admin only)
func (u *UseCase) ListTransactions(ctx context.Context) (*ListTransactionsResult, error) {
	transactions, err := u.transactionRepository.List(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch transactions")
	}

	return &ListTransactionsResult{
		Transactions: transactions,
	}, nil
}

// GetTransactionByID gets a transaction by ID
func (u *UseCase) GetTransactionByID(ctx context.Context, id uint) (*GetTransactionByIDResult, error) {
	txEntity, err := u.transactionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	return &GetTransactionByIDResult{
		Transaction: txEntity,
	}, nil
}

// UpdateTransactionStatus updates transaction status
func (u *UseCase) UpdateTransactionStatus(ctx context.Context, p *UpdateTransactionStatusParams) error {
	if err := u.transactionRepository.UpdateStatus(ctx, p.TransactionID, p.Status); err != nil {
		return errors.New("failed to update transaction status")
	}

	return nil
}

// ListTransactionsByBuyerID lists transactions by buyer ID
func (u *UseCase) ListTransactionsByBuyerID(ctx context.Context, buyerID uint) (*ListTransactionsResult, error) {
	transactions, err := u.transactionRepository.ListByBuyerID(ctx, buyerID)
	if err != nil {
		return nil, errors.New("failed to fetch your transactions")
	}

	return &ListTransactionsResult{
		Transactions: transactions,
	}, nil
}

// ListTransactionsBySellerID lists transactions by seller ID
func (u *UseCase) ListTransactionsBySellerID(ctx context.Context, sellerID uint) (*ListTransactionsResult, error) {
	transactions, err := u.transactionRepository.ListBySellerID(ctx, sellerID)
	if err != nil {
		return nil, errors.New("failed to fetch your transactions")
	}

	return &ListTransactionsResult{
		Transactions: transactions,
	}, nil
}
