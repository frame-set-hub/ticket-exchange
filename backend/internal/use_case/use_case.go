package use_case

import (
	"github.com/TicketX/backend/internal/entity/user"
	"github.com/TicketX/backend/internal/repository/message_repository"
	"github.com/TicketX/backend/internal/repository/ticket_repository"
	"github.com/TicketX/backend/internal/repository/transaction_repository"
	"github.com/TicketX/backend/internal/repository/user_repository"
)

// UseCase orchestrates application business logic
type UseCase struct {
	userRepository        user_repository.UserRepository
	ticketRepository      ticket_repository.TicketRepository
	transactionRepository transaction_repository.TransactionRepository
	messageRepository     message_repository.MessageRepository
}

// Config contains use case configuration
type Config struct {
	AdminUserID uint
}

// New creates a new UseCase instance
func New(
	config Config,
	userRepo user_repository.UserRepository,
	ticketRepo ticket_repository.TicketRepository,
	txRepo transaction_repository.TransactionRepository,
	msgRepo message_repository.MessageRepository,
) *UseCase {
	return &UseCase{
		userRepository:        userRepo,
		ticketRepository:      ticketRepo,
		transactionRepository: txRepo,
		messageRepository:     msgRepo,
	}
}

// GetUserRepository returns user repository (for testing)
func (u *UseCase) GetUserRepository() user_repository.UserRepository {
	return u.userRepository
}

// GetTicketRepository returns ticket repository (for testing)
func (u *UseCase) GetTicketRepository() ticket_repository.TicketRepository {
	return u.ticketRepository
}

// GetTransactionRepository returns transaction repository (for testing)
func (u *UseCase) GetTransactionRepository() transaction_repository.TransactionRepository {
	return u.transactionRepository
}

// GetMessageRepository returns message repository (for testing)
func (u *UseCase) GetMessageRepository() message_repository.MessageRepository {
	return u.messageRepository
}

// IsAdmin checks if user is admin
func (u *UseCase) IsAdmin(user *user.User) bool {
	return user.Role == "Admin"
}
