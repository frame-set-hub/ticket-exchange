package use_case

import (
	"context"
	"errors"

	"github.com/TicketX/backend/internal/entity/user"
	"github.com/TicketX/backend/pkg/utils"
)

// RegisterParams contains parameters for user registration
type RegisterParams struct {
	Username string
	Email    string
	Password string
	Role     string
}

// RegisterResult contains result of user registration
type RegisterResult struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// LoginParams contains parameters for user login
type LoginParams struct {
	Username string
	Password string
}

// LoginResult contains result of user login
type LoginResult struct {
	Token  string        `json:"token"`
	User   *user.User    `json:"user"`
}

// Register registers a new user
func (u *UseCase) Register(ctx context.Context, p *RegisterParams) (*RegisterResult, error) {
	// Validate role
	role := user.RoleUser
	if p.Role == string(user.RoleAdmin) {
		role = user.RoleAdmin
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(p.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	userEntity := &user.User{
		Username: p.Username,
		Email:    p.Email,
		Password: hashedPassword,
		Role:     role,
	}

	if err := u.userRepository.Create(ctx, userEntity); err != nil {
		if errors.Is(err, errors.New("username or email already exists")) {
			return nil, errors.New("username or email already exists")
		}
		return nil, err
	}

	return &RegisterResult{
		ID:       userEntity.ID,
		Username: userEntity.Username,
		Email:    userEntity.Email,
		Role:     string(userEntity.Role),
	}, nil
}

// Login authenticates a user
func (u *UseCase) Login(ctx context.Context, p *LoginParams) (*LoginResult, error) {
	// Find user by username
	userEntity, err := u.userRepository.FindByUsername(ctx, p.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	if !utils.CheckPasswordHash(p.Password, userEntity.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(userEntity.ID, userEntity.Username, userEntity.Email, userEntity.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &LoginResult{
		Token: token,
		User:  userEntity,
	}, nil
}

// GetUserByID gets a user by ID
func (u *UseCase) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	return u.userRepository.FindByID(ctx, id)
}
