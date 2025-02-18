package services

import (
	"context"
	"fmt"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

type KeycloakService struct {
	client       *gocloak.GoCloak
	token        *gocloak.JWT
	realm        string
	clientID     string
	clientSecret string
}

// NewKeycloakService initializes Keycloak client and logs in
func NewKeycloakService() (*KeycloakService, error) {
	client := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_URL"))
	realm := os.Getenv("KEYCLOAK_REALM")
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")

	// Authenticate with Keycloak using context
	ctx := context.TODO()
	token, err := client.LoginClient(ctx, clientID, clientSecret, realm)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate Keycloak client: %w", err)
	}

	return &KeycloakService{
		client:       client,
		token:        token,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
	}, nil
}

// GetUserDetails fetches user details from Keycloak
func (k *KeycloakService) GetUserDetails(userID string) (*gocloak.User, error) {
	ctx := context.TODO()
	user, err := k.client.GetUserByID(ctx, k.token.AccessToken, k.realm, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user details from Keycloak: %w", err)
	}
	return user, nil
}
