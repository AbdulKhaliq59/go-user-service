package middleware

import (
	"net/http"
	"strings"
	"time"

	"go-user-service/models"
	"go-user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(keycloakService *services.KeycloakService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, &models.Response{
				Timestamp: time.Now(),
				Message:   "Authorization header required",
				Status:    http.StatusUnauthorized,
				Data:      nil,
			})
			c.Abort()
			return
		}

		// Extract Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, &models.Response{
				Timestamp: time.Now(),
				Message:   "Invalid authorization format",
				Status:    http.StatusUnauthorized,
				Data:      nil,
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		claims := &services.KeycloakClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Get public key from Keycloak
			return keycloakService.GetPublicKey()
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, &models.Response{
				Timestamp: time.Now(),
				Message:   "Invalid or expired token",
				Status:    http.StatusUnauthorized,
				Data:      nil,
			})
			c.Abort()
			return
		}

		// Set user claims in context
		c.Set("user", *claims)
		c.Next()
	}
}

// HasRole checks if the authenticated user has the required role
func HasRole(c *gin.Context, requiredRole string) bool {
	userClaims, exists := c.Get("user")
	if !exists {
		return false
	}

	// Type assert to KeycloakClaims
	claims, ok := userClaims.(services.KeycloakClaims)
	if !ok {
		return false
	}

	// Check for required role in the claims
	// Note: You'll need to modify this based on how roles are stored in your JWT
	// This is a simplified example
	realm := claims.RegisteredClaims.Issuer
	if realm == "" {
		return false
	}

	// Check for the role in realm_access.roles claim
	// This would need to be adapted based on your actual JWT structure
	// For example, you might need to check a "realm_access" claim

	// For now, this is a placeholder implementation
	// In a real implementation, you'd extract roles from the token
	realmRoles := []string{"realm:assign_users_to_group", "realm:manage_users"}

	for _, role := range realmRoles {
		if role == requiredRole {
			return true
		}
	}

	return false
}
