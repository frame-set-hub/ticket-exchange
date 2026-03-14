package message_repository

import (
	"context"
	"time"

	"github.com/TicketX/backend/internal/entity/message"
	"gorm.io/gorm"
)

// pgDb implements MessageRepository interface with PostgreSQL
type pgDb struct {
	db *gorm.DB
}

// messageModel represents the database model for message
type messageModel struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	TransactionID uint
	SenderID      uint
	ReceiverID    uint
	Content       string
	AttachmentURL string
}

// NewPostgresDb creates a new PostgreSQL message repository
func NewPostgresDb(db *gorm.DB) MessageRepository {
	return &pgDb{db: db}
}

// Create creates a new message
func (r *pgDb) Create(ctx context.Context, m *message.Message) error {
	msgModel := &messageModel{
		TransactionID: m.TransactionID,
		SenderID:      m.SenderID,
		ReceiverID:    m.ReceiverID,
		Content:       m.Content,
		AttachmentURL: m.AttachmentURL,
	}

	if err := r.db.WithContext(ctx).Create(msgModel).Error; err != nil {
		return err
	}

	m.ID = msgModel.ID
	return nil
}

// ListByTransactionID lists messages by transaction ID
func (r *pgDb) ListByTransactionID(ctx context.Context, transactionID uint) ([]*message.Message, error) {
	var msgModels []messageModel

	if err := r.db.WithContext(ctx).
		Where("transaction_id = ?", transactionID).
		Order("created_at ASC").
		Find(&msgModels).Error; err != nil {
		return nil, err
	}

	result := make([]*message.Message, len(msgModels))
	for i, msg := range msgModels {
		result[i] = &message.Message{
			ID:            msg.ID,
			CreatedAt:     msg.CreatedAt,
			TransactionID: msg.TransactionID,
			SenderID:      msg.SenderID,
			ReceiverID:    msg.ReceiverID,
			Content:       msg.Content,
			AttachmentURL: msg.AttachmentURL,
		}
	}

	return result, nil
}
