package user_repository

import (
	"context"

	"github.com/TicketX/backend/internal/entity/user"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByID(ctx context.Context, id uint) (*user.User, error)
}
