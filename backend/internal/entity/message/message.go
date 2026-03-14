package message

import "time"

// Message represents the domain entity for chat message
type Message struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	TransactionID  uint      `json:"transaction_id"`
	SenderID       uint      `json:"sender_id"`
	SenderUsername string    `json:"sender_username,omitempty"`
	ReceiverID     uint      `json:"receiver_id"`
	Content        string    `json:"content"`
	AttachmentURL  string    `json:"attachment_url"`
}
