package user

import "time"

// User represents the domain entity for user
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // Include for login, but don't expose in JSON
	Role     Role   `json:"role"`
}

// UserWithCreated represents user with created_at field
type UserWithCreated struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Role represents user roles
type Role string

const (
	RoleUser  Role = "User"
	RoleAdmin Role = "Admin"
)
