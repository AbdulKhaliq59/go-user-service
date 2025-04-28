package middleware

import (
	"net/http"
	"time"

	"go-user-service/models"
	"go-user-service/services"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware creates a middleware for role-based access control
func RoleMiddleware(requiredRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user claims from context
		userClaims, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, &models.Response{
				Timestamp: time.Now(),
				Message:   "User authentication required",
				Status:    http.StatusUnauthorized,
				Data:      nil,
			})
			c.Abort()
			return
		}

		// Type assert to KeycloakClaims
		claims, ok := userClaims.(services.KeycloakClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, &models.Response{
				Timestamp: time.Now(),
				Message:   "Invalid user claims format",
				Status:    http.StatusInternalServerError,
				Data:      nil,
			})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		// Note: You'll need to implement the HasRole method in KeycloakClaims
		hasAccess := false
		for _, role := range requiredRoles {
			if claims.HasRole(role) {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, &models.Response{
				Timestamp: time.Now(),
				Message:   "Insufficient permissions",
				Status:    http.StatusForbidden,
				Data:      nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
