package message_repository

import (
	"context"

	"github.com/TicketX/backend/internal/entity/message"
)

// MessageRepository defines the interface for message data operations
type MessageRepository interface {
	Create(ctx context.Context, m *message.Message) error
	ListByTransactionID(ctx context.Context, transactionID uint) ([]*message.Message, error)
}
