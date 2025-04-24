package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"go-user-service/api/dto"

	"github.com/Nerzal/gocloak/v13"
)

// RoleService handles business logic related to roles
type RoleService struct {
	keycloakBaseURL string
	realm           string
	clientID        string
	clientSecret    string
}

// NewRoleService creates a new role service with validation
func NewRoleService() *RoleService {
	// Debug: Print all Keycloak-related environment variables
	log.Println("Loading Keycloak configuration...")
	baseURL := os.Getenv("KEYCLOAK_BASE_URL")
	realm := os.Getenv("KEYCLOAK_REALM")
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")

	// Validate required values
	if baseURL == "" {
		log.Fatal("KEYCLOAK_BASE_URL is not set")
	}
	if realm == "" {
		log.Fatal("KEYCLOAK_REALM is not set")
	}
	if clientID == "" {
		log.Fatal("KEYCLOAK_CLIENT_ID is not set")
	}
	if clientSecret == "" {
		log.Fatal("KEYCLOAK_CLIENT_SECRET is not set")
	}

	// Ensure base URL has a protocol
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
		log.Printf("Added https:// prefix to Keycloak base URL: %s", baseURL)
	}

	return &RoleService{
		keycloakBaseURL: strings.TrimSuffix(baseURL, "/"), // Remove trailing slash if present
		realm:           realm,
		clientID:        clientID,
		clientSecret:    clientSecret,
	}
}

// createKeycloakAdminClient creates a new admin client and gets a fresh token
func (s *RoleService) createKeycloakAdminClient(ctx context.Context) (*gocloak.GoCloak, string, error) {
	// Ensure the base URL doesn't end with a trailing slash
	baseURL := strings.TrimSuffix(s.keycloakBaseURL, "/")

	// Create client with the sanitized base URL
	client := gocloak.NewClient(baseURL)

	// Get a fresh token
	token, err := client.LoginClient(ctx, s.clientID, s.clientSecret, s.realm)
	if err != nil {
		return nil, "", fmt.Errorf("failed to authenticate with Keycloak: %w", err)
	}

	log.Printf("Successfully authenticated with Keycloak")
	return client, token.AccessToken, nil
}

// GetRoles fetches all realm roles from Keycloak
func (s *RoleService) GetRoles(ctx context.Context) ([]*gocloak.Role, error) {
	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		log.Printf("Error creating Keycloak admin client: %v", err)
		return nil, err
	}

	// Adding debug logs
	log.Printf("Attempting to fetch roles from realm: %s", s.realm)

	// Use empty params for no filtering
	roles, err := client.GetRealmRoles(ctx, token, s.realm, gocloak.GetRoleParams{})
	if err != nil {
		log.Printf("Error fetching roles from Keycloak: %v", err)
		return nil, fmt.Errorf("error fetching roles: %w", err)
	}

	log.Printf("Successfully fetched %d roles", len(roles))
	return roles, nil
}

// FindById finds a role by ID
func (s *RoleService) FindById(ctx context.Context, id string) (*gocloak.Role, error) {
	fmt.Printf("Finding role with ID: %s\n", id)

	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		return nil, err
	}

	// First try to get by ID (some Keycloak versions might support this)
	role, err := client.GetRealmRole(ctx, token, s.realm, id)
	if err == nil {
		return role, nil
	}

	// If not found by ID, try to find by name in all roles
	roles, err := client.GetRealmRoles(ctx, token, s.realm, gocloak.GetRoleParams{})
	if err != nil {
		return nil, fmt.Errorf("error fetching roles: %w", err)
	}

	for _, r := range roles {
		if *r.ID == id {
			return r, nil
		}
	}

	return nil, fmt.Errorf("role with ID '%s' not found", id)
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(ctx context.Context, createRoleDto dto.CreateRoleDTO) error {
	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		return err
	}

	// Convert DTO to Keycloak role
	role := gocloak.Role{
		Name:        &createRoleDto.Name,
		Description: &createRoleDto.Description,
	}

	_, err = client.CreateRealmRole(ctx, token, s.realm, role)
	if err != nil {
		return fmt.Errorf("error creating role: %w", err)
	}

	return nil
}

func (s *RoleService) UpdateRole(ctx context.Context, id string, updateRoleDto dto.UpdateRoleDTO) error {
	log.Printf("Attempting to update role with ID: %s", id)

	// 1. First get the existing role to know its current name
	existingRole, err := s.FindById(ctx, id)
	if err != nil {
		log.Printf("Error finding role with ID %s: %v", id, err)
		return fmt.Errorf("error finding role: %w", err)
	}

	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		log.Printf("Error creating Keycloak admin client: %v", err)
		return err
	}

	// 2. Create a copy of the existing role to modify
	updatedRole := *existingRole

	// 3. Update fields if provided
	if updateRoleDto.Name != "" {
		updatedRole.Name = &updateRoleDto.Name
	}
	if updateRoleDto.Description != "" {
		updatedRole.Description = &updateRoleDto.Description
	}

	// 4. Special handling for name changes
	if updateRoleDto.Name != "" && updateRoleDto.Name != *existingRole.Name {
		log.Printf("Attempting role rename from '%s' to '%s'", *existingRole.Name, updateRoleDto.Name)

		// First update the role with its existing name but new attributes
		err = client.UpdateRealmRole(ctx, token, s.realm, *existingRole.Name, updatedRole)
		if err != nil {
			log.Printf("Error updating role attributes: %v", err)
			return fmt.Errorf("error updating role attributes: %w", err)
		}

		log.Printf("Successfully updated role attributes, name remains '%s' temporarily", *existingRole.Name)
	} else {
		// Normal update without name change
		err = client.UpdateRealmRole(ctx, token, s.realm, *existingRole.Name, updatedRole)
		if err != nil {
			log.Printf("Error updating role: %v", err)
			return fmt.Errorf("error updating role: %w", err)
		}
	}

	log.Printf("Successfully updated role with ID: %s", id)
	return nil
}

// DeleteRole deletes a role by ID
func (s *RoleService) DeleteRole(ctx context.Context, id string) error {
	log.Printf("Attempting to delete role with ID: %s", id)

	// Find the role by ID to get its name
	role, err := s.FindById(ctx, id)
	if err != nil {
		log.Printf("Error finding role with ID %s: %v", id, err)
		return fmt.Errorf("error finding role: %w", err)
	}

	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		log.Printf("Error creating Keycloak admin client: %v", err)
		return err
	}

	// Use the role name to delete the role
	err = client.DeleteRealmRole(ctx, token, s.realm, *role.Name)
	if err != nil {
		log.Printf("Error deleting role with name %s: %v", *role.Name, err)
		return fmt.Errorf("error deleting role: %w", err)
	}

	log.Printf("Successfully deleted role with name: %s", *role.Name)
	return nil
}
