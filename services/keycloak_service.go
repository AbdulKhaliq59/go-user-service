package services

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
)

type KeycloakService struct {
	client       *gocloak.GoCloak
	token        *gocloak.JWT
	realm        string
	clientID     string
	clientSecret string
	baseURL      string
	groupService *GroupService
}

type KeycloakClaims struct {
	jwt.RegisteredClaims
	EmailVerified  bool             `json:"email_verified"`
	PreferredName  string           `json:"preferred_username"`
	GivenName      string           `json:"given_name"`
	FamilyName     string           `json:"family_name"`
	Email          string           `json:"email"`
	RealmAccess    RealmAccess      `json:"realm_access"`
	ResourceAccess map[string]Roles `json:"resource_access"`
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type Roles struct {
	Roles []string `json:"roles"`
}

// NewKeycloakService initializes Keycloak client and logs in
func NewKeycloakService() (*KeycloakService, error) {
	baseURL := os.Getenv("KEYCLOAK_BASE_URL")
	client := gocloak.NewClient(baseURL)
	realm := os.Getenv("KEYCLOAK_REALM")
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")

	// Authenticate with Keycloak using context
	ctx := context.TODO()
	token, err := client.LoginClient(ctx, clientID, clientSecret, realm)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate Keycloak client: %w", err)
	}

	groupService := NewGroupService()

	return &KeycloakService{
		client:       client,
		token:        token,
		realm:        realm,
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      baseURL,
		groupService: groupService,
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

// RefreshAdminToken refreshes the admin token if needed
func (k *KeycloakService) RefreshAdminToken() error {
	ctx := context.TODO()
	token, err := k.client.LoginClient(ctx, k.clientID, k.clientSecret, k.realm)
	if err != nil {
		return fmt.Errorf("failed to refresh Keycloak admin token: %w", err)
	}
	k.token = token
	return nil
}

// GetPublicKey fetches the public key from Keycloak
func (k *KeycloakService) GetPublicKey() (interface{}, error) {
	return nil, errors.New("Method not implemented")
}

func (kc *KeycloakClaims) HasRole(role string) bool {
	// Check realm roles
	for _, r := range kc.RealmAccess.Roles {
		if r == role {
			return true
		}
	}

	// Check client-specific roles
	for _, roles := range kc.ResourceAccess {
		for _, r := range roles.Roles {
			if r == role {
				return true
			}
		}
	}

	return false
}
