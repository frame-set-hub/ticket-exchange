package use_case

import (
	"context"
	"errors"

	"github.com/TicketX/backend/internal/entity/message"
)

// CreateMessageParams contains parameters for creating a message
type CreateMessageParams struct {
	TransactionID uint
	SenderID      uint
	ReceiverID    uint
	Content       string
	AttachmentURL string
}

// CreateMessageResult contains result of creating a message
type CreateMessageResult struct {
	Message *message.Message
}

// ListMessagesByTransactionIDResult contains result of listing messages
type ListMessagesByTransactionIDResult struct {
	Messages []*message.Message
}

// CreateMessage creates a new message
func (u *UseCase) CreateMessage(ctx context.Context, p *CreateMessageParams) (*CreateMessageResult, error) {
	msgEntity := &message.Message{
		TransactionID: p.TransactionID,
		SenderID:      p.SenderID,
		ReceiverID:    p.ReceiverID,
		Content:       p.Content,
		AttachmentURL: p.AttachmentURL,
	}

	if err := u.messageRepository.Create(ctx, msgEntity); err != nil {
		return nil, errors.New("failed to create message")
	}

	return &CreateMessageResult{
		Message: msgEntity,
	}, nil
}

// ListMessagesByTransactionID lists messages for a transaction
func (u *UseCase) ListMessagesByTransactionID(ctx context.Context, transactionID uint) (*ListMessagesByTransactionIDResult, error) {
	messages, err := u.messageRepository.ListByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, errors.New("failed to fetch messages")
	}

	// Enrich with sender usernames
	usernameCache := make(map[uint]string)
	for _, msg := range messages {
		if _, ok := usernameCache[msg.SenderID]; !ok {
			if user, err := u.userRepository.FindByID(ctx, msg.SenderID); err == nil {
				usernameCache[msg.SenderID] = user.Username
			}
		}
		msg.SenderUsername = usernameCache[msg.SenderID]
	}

	return &ListMessagesByTransactionIDResult{
		Messages: messages,
	}, nil
}
