package services

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware is a middleware for authenticating requests using JWT
type AuthMiddleware struct {
	logger      *logrus.Logger
	authService *AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(logger *logrus.Logger, authService *AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		logger:      logger,
		authService: authService,
	}
}

// JWT returns a middleware function that validates JWT tokens
func (m *AuthMiddleware) JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Fallback to token from query parameter
			authHeader = c.Query("token")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"code": 4001, "msg": "Authorization token required"})
				c.Abort()
				return
			}
		} else {
			// Split Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				c.JSON(http.StatusUnauthorized, gin.H{"code": 4001, "msg": "Authorization header format must be Bearer {token}"})
				c.Abort()
				return
			}
			authHeader = parts[1]
		}

		// Parse and validate token
		claims, err := m.authService.ValidateToken(authHeader)
		if err != nil {
			m.logger.WithError(err).Warn("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"code": 4001, "msg": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID and username in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
