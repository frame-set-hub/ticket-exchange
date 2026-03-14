package user_repository

import (
	"context"
	"errors"

	"time"

	"github.com/TicketX/backend/internal/entity/user"
	"gorm.io/gorm"
)

// pgDb implements UserRepository interface with PostgreSQL
type pgDb struct {
	db *gorm.DB
}

// userModel represents the database model for user
type userModel struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string         `gorm:"uniqueIndex;not null"`
	Email     string         `gorm:"uniqueIndex;not null"`
	Password  string         `gorm:"not null"`
	Role      user.Role      `gorm:"type:varchar(20);default:'User'"`
}

// NewPostgresDb creates a new PostgreSQL user repository
func NewPostgresDb(db *gorm.DB) UserRepository {
	return &pgDb{db: db}
}

// Create creates a new user
func (r *pgDb) Create(ctx context.Context, u *user.User) error {
	userModel := &userModel{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
	}

	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("username or email already exists")
		}
		return err
	}

	u.ID = userModel.ID
	return nil
}

// FindByUsername finds a user by username (includes password for login)
func (r *pgDb) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var userModel userModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user.User{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Password: userModel.Password, // Include password for login verification
		Role:     userModel.Role,
	}, nil
}

// FindByEmail finds a user by email
func (r *pgDb) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var userModel userModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user.User{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Password: userModel.Password,
		Role:     userModel.Role,
	}, nil
}

// FindByID finds a user by ID
func (r *pgDb) FindByID(ctx context.Context, id uint) (*user.User, error) {
	var userModel userModel
	if err := r.db.WithContext(ctx).First(&userModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user.User{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Role:     userModel.Role,
	}, nil
}
