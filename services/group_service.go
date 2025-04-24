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

// GroupService handles business logic related to groups
type GroupService struct {
	keycloakBaseURL string
	realm           string
	clientID        string
	clientSecret    string
}

// NewGroupService creates a new group service with validation
func NewGroupService() *GroupService {
	// Debug: Print all Keycloak-related environment variables
	log.Println("Loading Keycloak configuration for group service...")
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

	return &GroupService{
		keycloakBaseURL: strings.TrimSuffix(baseURL, "/"), // Remove trailing slash if present
		realm:           realm,
		clientID:        clientID,
		clientSecret:    clientSecret,
	}
}

// createKeycloakAdminClient creates a new admin client and gets a fresh token
func (s *GroupService) createKeycloakAdminClient(ctx context.Context) (*gocloak.GoCloak, string, error) {
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

// GetGroups fetches all groups from Keycloak
func (s *GroupService) GetGroups(ctx context.Context) ([]*gocloak.Group, error) {
	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		log.Printf("Error creating Keycloak admin client: %v", err)
		return nil, err
	}

	// Adding debug logs
	log.Printf("Attempting to fetch groups from realm: %s", s.realm)

	// Get all groups
	groups, err := client.GetGroups(ctx, token, s.realm, gocloak.GetGroupsParams{})
	if err != nil {
		log.Printf("Error fetching groups from Keycloak: %v", err)
		return nil, fmt.Errorf("error fetching groups: %w", err)
	}

	log.Printf("Successfully fetched %d groups", len(groups))
	return groups, nil
}

// GetGroupById finds a group by ID
func (s *GroupService) GetGroupById(ctx context.Context, id string) (*gocloak.Group, error) {
	fmt.Printf("Finding group with ID: %s\n", id)

	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		return nil, err
	}

	group, err := client.GetGroup(ctx, token, s.realm, id)
	if err != nil {
		return nil, fmt.Errorf("error finding group with ID '%s': %w", id, err)
	}

	return group, nil
}

// GetGroupRoles fetches all roles assigned to a group
func (s *GroupService) GetGroupRoles(ctx context.Context, groupID string) ([]*gocloak.Role, error) {
	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		return nil, err
	}

	// Get all realm roles assigned to the group
	roles, err := client.GetRealmRolesByGroupID(ctx, token, s.realm, groupID)
	if err != nil {
		return nil, fmt.Errorf("error fetching roles for group ID '%s': %w", groupID, err)
	}

	return roles, nil
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(ctx context.Context, createGroupDto dto.CreateGroupDTO) error {
	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		return err
	}

	// Convert DTO to Keycloak group
	group := gocloak.Group{
		Name: &createGroupDto.Name,
	}

	_, err = client.CreateGroup(ctx, token, s.realm, group)
	if err != nil {
		return fmt.Errorf("error creating group: %w", err)
	}

	return nil
}

// UpdateGroup updates an existing group
func (s *GroupService) UpdateGroup(ctx context.Context, id string, updateGroupDto dto.UpdateGroupDTO) error {
	log.Printf("Attempting to update group with ID: %s", id)

	// First get the existing group
	existingGroup, err := s.GetGroupById(ctx, id)
	if err != nil {
		log.Printf("Error finding group with ID %s: %v", id, err)
		return fmt.Errorf("error finding group: %w", err)
	}

	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		log.Printf("Error creating Keycloak admin client: %v", err)
		return err
	}

	// Create a copy of the existing group to modify
	updatedGroup := *existingGroup

	// Update fields if provided
	if updateGroupDto.Name != "" {
		updatedGroup.Name = &updateGroupDto.Name
	}

	// Update the group
	err = client.UpdateGroup(ctx, token, s.realm, updatedGroup)
	if err != nil {
		log.Printf("Error updating group: %v", err)
		return fmt.Errorf("error updating group: %w", err)
	}

	log.Printf("Successfully updated group with ID: %s", id)
	return nil
}

// DeleteGroup deletes a group by ID
func (s *GroupService) DeleteGroup(ctx context.Context, id string) error {
	log.Printf("Attempting to delete group with ID: %s", id)

	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		log.Printf("Error creating Keycloak admin client: %v", err)
		return err
	}

	// Delete the group
	err = client.DeleteGroup(ctx, token, s.realm, id)
	if err != nil {
		log.Printf("Error deleting group with ID %s: %v", id, err)
		return fmt.Errorf("error deleting group: %w", err)
	}

	log.Printf("Successfully deleted group with ID: %s", id)
	return nil
}

// AssignRolesToGroup assigns roles to a group
func (s *GroupService) AssignRolesToGroup(ctx context.Context, assignRolesDto dto.AssignRolesToGroupDTO) error {
	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		return err
	}

	// Find the group
	group, err := s.GetGroupById(ctx, assignRolesDto.GroupID)
	if err != nil {
		return fmt.Errorf("error finding group: %w", err)
	}

	// Prepare roles to assign
	var rolesToAssign []gocloak.Role
	for _, roleID := range assignRolesDto.RoleIDs {
		// Get role by ID
		role, err := client.GetRealmRole(ctx, token, s.realm, roleID)
		if err != nil {
			// Try to get role by name if ID doesn't work
			roles, err := client.GetRealmRoles(ctx, token, s.realm, gocloak.GetRoleParams{})
			if err != nil {
				return fmt.Errorf("error fetching roles: %w", err)
			}

			found := false
			for _, r := range roles {
				if *r.ID == roleID {
					rolesToAssign = append(rolesToAssign, *r)
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("role with ID '%s' not found", roleID)
			}
		} else {
			rolesToAssign = append(rolesToAssign, *role)
		}
	}

	// Assign roles to the group
	for _, role := range rolesToAssign {
		err = client.AddRealmRoleToGroup(ctx, token, s.realm, *group.ID, []gocloak.Role{role})
		if err != nil {
			return fmt.Errorf("error assigning role '%s' to group: %w", *role.Name, err)
		}
	}

	return nil
}

// UnassignRoleFromGroup removes a role from a group
func (s *GroupService) UnassignRoleFromGroup(ctx context.Context, unassignRoleDto dto.UnassignRoleFromGroupDTO) error {
	client, token, err := s.createKeycloakAdminClient(ctx)
	if err != nil {
		return err
	}

	// Find the group
	group, err := s.GetGroupById(ctx, unassignRoleDto.GroupID)
	if err != nil {
		return fmt.Errorf("error finding group: %w", err)
	}

	// Prepare roles to unassign
	var rolesToUnassign []gocloak.Role
	for _, roleID := range unassignRoleDto.RoleIDs {
		// Get role by ID
		role, err := client.GetRealmRole(ctx, token, s.realm, roleID)
		if err != nil {
			// Try to get role by name if ID doesn't work
			roles, err := client.GetRealmRoles(ctx, token, s.realm, gocloak.GetRoleParams{})
			if err != nil {
				return fmt.Errorf("error fetching roles: %w", err)
			}

			found := false
			for _, r := range roles {
				if *r.ID == roleID {
					rolesToUnassign = append(rolesToUnassign, *r)
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("role with ID '%s' not found", roleID)
			}
		} else {
			rolesToUnassign = append(rolesToUnassign, *role)
		}
	}

	// Get group members
	members, err := client.GetGroupMembers(ctx, token, s.realm, *group.ID, gocloak.GetGroupsParams{})
	if err != nil {
		return fmt.Errorf("error fetching group members: %w", err)
	}

	// Remove roles from each group member
	for _, member := range members {
		for _, role := range rolesToUnassign {
			err = client.DeleteRealmRoleFromUser(ctx, token, s.realm, *member.ID, []gocloak.Role{role})
			if err != nil {
				log.Printf("Warning: Failed to remove role '%s' from user '%s': %v", *role.Name, *member.Username, err)
			}
		}
	}

	// Remove roles from the group
	for _, role := range rolesToUnassign {
		err = client.DeleteRealmRoleFromGroup(ctx, token, s.realm, *group.ID, []gocloak.Role{role})
		if err != nil {
			return fmt.Errorf("error removing role '%s' from group: %w", *role.Name, err)
		}
	}

	return nil
}

// GetGroupsAndRoles fetches all groups and their assigned roles
func (s *GroupService) GetGroupsAndRoles(ctx context.Context) (map[string][]*gocloak.Role, error) {
	groups, err := s.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	groupRolesMap := make(map[string][]*gocloak.Role)

	// For each group, get its roles
	for _, group := range groups {
		roles, err := s.GetGroupRoles(ctx, *group.ID)
		if err != nil {
			log.Printf("Warning: Failed to fetch roles for group '%s': %v", *group.Name, err)
			groupRolesMap[*group.Name] = []*gocloak.Role{}
			continue
		}
		groupRolesMap[*group.Name] = roles
	}

	return groupRolesMap, nil
}
