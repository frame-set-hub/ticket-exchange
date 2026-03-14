package middleware

import (
	"net/http"
	"strings"

	"github.com/TicketX/backend/internal/entity/user"
	"github.com/TicketX/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthContextKey is the key used to store user info in context
type AuthContextKey string

const UserContextKey AuthContextKey = "user"

// Authenticate middleware validates JWT token and sets user in context
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check for Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user in context
		userInfo := &user.User{
			ID:       claims.UserID,
			Username: claims.Username,
			Email:    claims.Email,
			Role:     user.Role(claims.Role),
		}

		c.Set(string(UserContextKey), userInfo)
		c.Next()
	}
}

// AdminOnly middleware checks if user is admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, exists := c.Get(string(UserContextKey))
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user := userInfo.(*user.User)
		if user.Role != "Admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUser retrieves current user from context
func GetCurrentUser(c *gin.Context) *user.User {
	userInfo, exists := c.Get(string(UserContextKey))
	if !exists {
		return nil
	}
	return userInfo.(*user.User)
}
