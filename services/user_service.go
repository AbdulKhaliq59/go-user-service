// services/user_service.go
package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go-user-service/api/dto"
	"go-user-service/models"

	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	keycloakService *KeycloakService
}

func NewAuthService(keycloakService *KeycloakService) *AuthService {
	return &AuthService{
		keycloakService: keycloakService,
	}
}

type KeycloakClaims struct {
	jwt.RegisteredClaims
	EmailVerified bool   `json:"email_verified"`
	PreferredName string `json:"preferred_username"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Email         string `json:"email"`
}

func (s *AuthService) Login(username, password string) (*models.Response, error) {
	ctx := context.Background()
	client := s.keycloakService.client
	realm := s.keycloakService.realm
	clientID := s.keycloakService.clientID
	clientSecret := s.keycloakService.clientSecret

	// Find user by username or email
	users, err := client.GetUsers(ctx, s.keycloakService.token.AccessToken, realm, gocloak.GetUsersParams{
		Username: &username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if len(users) == 0 {
		// Try to find by email
		users, err = client.GetUsers(ctx, s.keycloakService.token.AccessToken, realm, gocloak.GetUsersParams{
			Email: &username,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to find user by email: %w", err)
		}
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := users[0]

	// Check if password update is required
	requiresPasswordUpdate := false
	if user.RequiredActions != nil {
		for _, action := range *user.RequiredActions {
			if action == "UPDATE_PASSWORD" {
				requiresPasswordUpdate = true
				break
			}
		}
	}

	if requiresPasswordUpdate {
		return &models.Response{
			Timestamp: time.Now(),
			Message:   "Password update required",
			Status:    http.StatusOK,
			Data:      nil,
		}, nil
	}

	// Get user token
	tokenResponse, err := client.GetToken(ctx, realm, gocloak.TokenOptions{
		ClientID:     &clientID,
		ClientSecret: &clientSecret,
		Username:     user.Username,
		Password:     &password,
		GrantType:    gocloak.StringP("password"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Create the response
	response := &models.Response{
		Timestamp: time.Now(),
		Message:   "Successful",
		Status:    http.StatusOK,
		Data: models.LoginResponseData{
			Token: tokenResponse.AccessToken,
		},
	}

	return response, nil
}

// Add to services/user_service.go

func (s *AuthService) CreateNewUser(createUserDto *dto.CreateUserDto) error {
	ctx := context.Background()
	client := s.keycloakService.client
	realm := s.keycloakService.realm

	// Create user representation for Keycloak
	enabled := true
	temporary := false
	newUser := gocloak.User{
		Username:  &createUserDto.Username,
		Enabled:   &enabled,
		Email:     &createUserDto.Email,
		LastName:  &createUserDto.LastName,
		FirstName: &createUserDto.FirstName,
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Type:      gocloak.StringP("password"),
				Value:     &createUserDto.Password,
				Temporary: &temporary,
			},
		},
		Attributes: &map[string][]string{
			"phoneNumber": {createUserDto.PhoneNumber},
		},
	}

	// Create the user in Keycloak
	_, err := client.CreateUser(ctx, s.keycloakService.token.AccessToken, realm, newUser)
	if err != nil {
		// Handle specific error cases
		if err.Error() == "409 Conflict" || err.Error() == "Conflict" {
			if containsSubstring(err.Error(), "username") {
				return fmt.Errorf("username is already taken")
			} else if containsSubstring(err.Error(), "email") {
				return fmt.Errorf("email is already taken")
			}
		}
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func (s *AuthService) Logout(refreshToken string) error {
	// Make sure we have a valid admin token
	if err := s.keycloakService.RefreshAdminToken(); err != nil {
		return fmt.Errorf("failed to refresh admin token: %w", err)
	}

	keycloakBaseURL := s.keycloakService.baseURL
	realm := s.keycloakService.realm
	clientID := s.keycloakService.clientID
	clientSecret := s.keycloakService.clientSecret

	// Create the form data
	formData := url.Values{}
	formData.Add("client_id", clientID)
	formData.Add("client_secret", clientSecret)
	formData.Add("refresh_token", refreshToken)

	// Build the logout URL
	logoutURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", keycloakBaseURL, realm)

	// Create the request
	req, err := http.NewRequest("POST", logoutURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}

	// Set the headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute logout request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("logout failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (s *AuthService) Create(createUserDto *dto.UserDto) error {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Create user representation for Keycloak
	enabled := true
	temporary := false
	newUser := gocloak.User{
		Username:  &createUserDto.Username,
		Enabled:   &enabled,
		Email:     &createUserDto.Email,
		LastName:  &createUserDto.LastName,
		FirstName: &createUserDto.FirstName,
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Type:      gocloak.StringP("password"),
				Value:     &createUserDto.Password,
				Temporary: &temporary,
			},
		},
		Attributes: &map[string][]string{
			"phoneNumber": {createUserDto.PhoneNumber},
		},
	}

	// Create the user in Keycloak
	_, err = client.CreateUser(ctx, token, realm, newUser)
	if err != nil {
		// Handle specific error cases
		if err.Error() == "409 Conflict" || err.Error() == "Conflict" {
			if containsSubstring(err.Error(), "username") {
				return fmt.Errorf("username is already taken")
			} else if containsSubstring(err.Error(), "email") {
				return fmt.Errorf("email is already taken")
			}
		}
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// CreatePlayer creates a new player user and assigns player group and roles
func (s *AuthService) CreatePlayer(createUserDto *dto.CreateUserDto) error {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Create user representation for Keycloak
	enabled := true
	temporary := false
	newUser := gocloak.User{
		Username:  &createUserDto.Username,
		Enabled:   &enabled,
		Email:     &createUserDto.Email,
		LastName:  &createUserDto.LastName,
		FirstName: &createUserDto.FirstName,
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Type:      gocloak.StringP("password"),
				Value:     &createUserDto.Password,
				Temporary: &temporary,
			},
		},
		Attributes: &map[string][]string{
			"phoneNumber": {createUserDto.PhoneNumber},
		},
	}

	// Create the user in Keycloak
	userID, err := client.CreateUser(ctx, token, realm, newUser)
	if err != nil {
		// Handle specific error cases
		if err.Error() == "409 Conflict" || err.Error() == "Conflict" {
			if containsSubstring(err.Error(), "username") {
				return fmt.Errorf("username is already taken")
			} else if containsSubstring(err.Error(), "email") {
				return fmt.Errorf("email is already taken")
			}
		}
		return fmt.Errorf("error creating user: %w", err)
	}

	// Find player group by ID
	// Using the hardcoded ID from the NestJS code
	playerGroupID := "12ce86cf-b338-4726-915b-30c3dd393054"
	group, err := s.keycloakService.groupService.FindById(playerGroupID)
	if err != nil {
		return fmt.Errorf("failed to find player group: %w", err)
	}

	// Add user to player group
	err = client.AddUserToGroup(ctx, token, realm, userID, *group.ID)
	if err != nil {
		return fmt.Errorf("failed to add user to player group: %w", err)
	}

	// Get group roles
	groupRoles, err := s.keycloakService.groupService.GetGroupRoles(playerGroupID)
	if err != nil {
		return fmt.Errorf("failed to get group roles: %w", err)
	}

	// Convert group roles to role mappings
	var roles []gocloak.Role
	for _, role := range groupRoles {
		if role.Name != nil && role.ID != nil {
			roles = append(roles, *role)
		}
	}

	// Assign group roles to the user
	if len(roles) > 0 {
		err = client.AddRealmRoleToUser(ctx, token, realm, userID, roles)
		if err != nil {
			return fmt.Errorf("failed to assign roles to user: %w", err)
		}
	}

	return nil
}

// Helper function to extract role names from []gocloak.Role
func extractRoleNames(roles []gocloak.Role) []string {
	roleNames := make([]string, 0, len(roles))
	for _, role := range roles {
		if role.Name != nil {
			roleNames = append(roleNames, *role.Name)
		}
	}
	return roleNames
}

// Helper function to convert []gocloak.Role to []*gocloak.Role

// GetUserInfo fetches detailed information about a user
func (s *AuthService) GetUserInfo(userId string) (map[string]interface{}, error) {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Get user details
	user, err := client.GetUserByID(ctx, token, realm, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Get user's realm roles
	realmRoles, err := client.GetRealmRolesByUserID(ctx, token, realm, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user realm roles: %w", err)
	}

	// Get user's groups
	groups, err := client.GetUserGroups(ctx, token, realm, userId, gocloak.GetGroupsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	// Extract role names from realm roles
	roleNames := make([]string, 0, len(realmRoles))
	for _, role := range realmRoles {
		if role.Name != nil {
			roleNames = append(roleNames, *role.Name)
		}
	}

	// Collect roles from groups
	var groupRoles []gocloak.Role
	for _, group := range groups {
		if group.ID == nil {
			continue
		}

		// Get roles for this group
		roles, err := client.GetRealmRolesByGroupID(ctx, token, realm, *group.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get roles for group %s: %w", *group.ID, err)
		}

		// Convert []*gocloak.Role to []gocloak.Role
		for _, role := range roles {
			if role != nil {
				groupRoles = append(groupRoles, *role)
			}
		}
	}

	// Create the response structure
	userInfo := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"firstName":  user.FirstName,
		"lastName":   user.LastName,
		"attributes": user.Attributes,
		"realmRoles": roleNames,
		"groups":     groups,
		"roles":      append(roleNames, extractRoleNames(groupRoles)...),
	}

	return userInfo, nil
}

// AssignGroupToUser assigns a user to a group and their associated roles
func (s *AuthService) AssignGroupToUser(assignDto *dto.AssignUserToGroupDto) error {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Verify user exists
	user, err := client.GetUserByID(ctx, token, realm, assignDto.UserId)
	if err != nil {
		return fmt.Errorf("user not found with id: %s", assignDto.UserId)
	}

	// Verify group exists
	group, err := client.GetGroup(ctx, token, realm, assignDto.GroupId)
	if err != nil {
		return fmt.Errorf("group not found with id: %s", assignDto.GroupId)
	}

	// Assign the group to the user
	err = client.AddUserToGroup(ctx, token, realm, *user.ID, *group.ID)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}

	// Get group roles
	groupRoles, err := client.GetRealmRolesByGroupID(ctx, token, realm, *group.ID)
	if err != nil {
		return fmt.Errorf("failed to get group roles: %w", err)
	}

	// Assign group roles to the user if there are any
	if len(groupRoles) > 0 {
		// Convert []*gocloak.Role to []gocloak.Role
		convertedRoles := make([]gocloak.Role, len(groupRoles))
		for i, role := range groupRoles {
			if role != nil {
				convertedRoles[i] = *role
			}
		}

		err = client.AddRealmRoleToUser(ctx, token, realm, *user.ID, convertedRoles)
		if err != nil {
			return fmt.Errorf("failed to assign roles to user: %w", err)
		}
	}

	return nil
}

func (s *AuthService) UnassignGroupFromUser(unassignDto *dto.UnassignUserFromGroupDto) error {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Verify user exists
	user, err := client.GetUserByID(ctx, token, realm, unassignDto.UserId)
	if err != nil {
		return fmt.Errorf("user not found with id: %s", unassignDto.UserId)
	}

	// Verify group exists
	group, err := client.GetGroup(ctx, token, realm, unassignDto.GroupId)
	if err != nil {
		return fmt.Errorf("group not found with id: %s", unassignDto.GroupId)
	}

	// Remove user from the group
	err = client.DeleteUserFromGroup(ctx, token, realm, *user.ID, *group.ID)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}

	// Get group roles
	groupRoles, err := client.GetRealmRolesByGroupID(ctx, token, realm, *group.ID)
	if err != nil {
		return fmt.Errorf("failed to get group roles: %w", err)
	}

	// Remove group roles from the user if there are any
	if len(groupRoles) > 0 {
		// Convert []*gocloak.Role to []gocloak.Role
		convertedRoles := make([]gocloak.Role, len(groupRoles))
		for i, role := range groupRoles {
			if role != nil {
				convertedRoles[i] = *role
			}
		}

		err = client.DeleteRealmRoleFromUser(ctx, token, realm, *user.ID, convertedRoles)
		if err != nil {
			return fmt.Errorf("failed to revoke roles from user: %w", err)
		}
	}

	return nil
}

// UnassignGroupsFromUser removes a user from multiple groups and revokes associated roles
func (s *AuthService) UnassignGroupsFromUser(unassignDto *dto.UnassignUserFromGroupsDto) error {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Verify user exists
	user, err := client.GetUserByID(ctx, token, realm, unassignDto.UserId)
	if err != nil {
		return fmt.Errorf("user not found with id: %s", unassignDto.UserId)
	}

	for _, groupId := range unassignDto.GroupIds {
		// Verify group exists
		group, err := client.GetGroup(ctx, token, realm, groupId)
		if err != nil {
			return fmt.Errorf("group not found with id: %s", groupId)
		}

		// Remove user from the group
		err = client.DeleteUserFromGroup(ctx, token, realm, *user.ID, *group.ID)
		if err != nil {
			return fmt.Errorf("failed to remove user from group %s: %w", groupId, err)
		}

		// Get group roles
		groupRoles, err := client.GetRealmRolesByGroupID(ctx, token, realm, *group.ID)
		if err != nil {
			return fmt.Errorf("failed to get group roles: %w", err)
		}

		// Remove group roles from the user if there are any
		if len(groupRoles) > 0 {
			// Convert []*gocloak.Role to []gocloak.Role
			convertedRoles := make([]gocloak.Role, len(groupRoles))
			for i, role := range groupRoles {
				if role != nil {
					convertedRoles[i] = *role
				}
			}

			err = client.DeleteRealmRoleFromUser(ctx, token, realm, *user.ID, convertedRoles)
			if err != nil {
				return fmt.Errorf("failed to revoke roles from user: %w", err)
			}
		}
	}

	return nil
}

func (s *AuthService) GetAllUsersWithGroupsAndRoles() ([]map[string]interface{}, error) {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Get all users
	users, err := client.GetUsers(ctx, token, realm, gocloak.GetUsersParams{
		Max: gocloak.IntP(1000), // Set a reasonable limit, adjust as needed
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	result := make([]map[string]interface{}, 0, len(users))

	// Process each user to get their groups and roles
	for _, user := range users {
		if user.ID == nil {
			continue
		}

		userId := *user.ID

		// Get user's groups
		groups, err := client.GetUserGroups(ctx, token, realm, userId, gocloak.GetGroupsParams{})
		if err != nil {
			return nil, fmt.Errorf("failed to get groups for user %s: %w", userId, err)
		}

		// Collect roles from all groups
		var roles []gocloak.Role
		for _, group := range groups {
			if group.ID == nil {
				continue
			}

			// Get roles for this group
			groupRoles, err := client.GetRealmRolesByGroupID(ctx, token, realm, *group.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get roles for group %s: %w", *group.ID, err)
			}

			// Convert []*gocloak.Role to []gocloak.Role and add to roles
			for _, role := range groupRoles {
				if role != nil {
					roles = append(roles, *role)
				}
			}
		}

		// Get user's direct realm roles
		realmRoles, err := client.GetRealmRolesByUserID(ctx, token, realm, userId)
		if err != nil {
			return nil, fmt.Errorf("failed to get realm roles for user %s: %w", userId, err)
		}

		// Add direct realm roles to the roles collection
		for _, role := range realmRoles {
			roles = append(roles, *role)
		}

		// Create a map for this user with all their information
		userData := map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"firstName":  user.FirstName,
			"lastName":   user.LastName,
			"attributes": user.Attributes,
			"enabled":    user.Enabled,
			"groups":     groups,
			"roles":      roles,
		}

		result = append(result, userData)
	}

	return result, nil
}

// UpdateUser updates an existing user's information in Keycloak
func (s *AuthService) UpdateUser(userId string, updatedUser *dto.UpdateUserDto) error {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Create user representation for Keycloak update
	user := gocloak.User{
		Username:  &updatedUser.Username,
		Email:     &updatedUser.Email,
		LastName:  &updatedUser.LastName,
		FirstName: &updatedUser.FirstName,
		Attributes: &map[string][]string{
			"phoneNumber": {updatedUser.PhoneNumber},
		},
	}

	// Update the user in Keycloak
	err = client.UpdateUser(ctx, token, realm, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// FindUserById fetches a user by their ID including roles
func (s *AuthService) FindUserById(userId string) (map[string]interface{}, error) {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Get user details
	user, err := client.GetUserByID(ctx, token, realm, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Get user's realm roles
	realmRoles, err := client.GetRealmRolesByUserID(ctx, token, realm, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user realm roles: %w", err)
	}

	// Extract role names from realm roles
	roleNames := make([]string, 0, len(realmRoles))
	for _, role := range realmRoles {
		if role.Name != nil {
			roleNames = append(roleNames, *role.Name)
		}
	}

	// Create the response structure
	userInfo := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"firstName":  user.FirstName,
		"lastName":   user.LastName,
		"attributes": user.Attributes,
		"realmRoles": roleNames,
	}

	return userInfo, nil
}

// ResetPassword resets a user's password
func (s *AuthService) ResetPassword(userId string, newPassword string) error {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Set password is temporary
	temporary := false

	// Create credential representation
	credential := gocloak.CredentialRepresentation{
		Type:      gocloak.StringP("password"),
		Value:     &newPassword,
		Temporary: &temporary,
	}

	// Reset the user's password
	err = client.SetPassword(ctx, token, realm, userId, *credential.Value, *credential.Temporary)
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}

// DeleteUserById deletes a user by their ID
func (s *AuthService) DeleteUserById(userId string) error {
	ctx := context.Background()

	// Create Keycloak admin client through GroupService
	client, token, err := s.keycloakService.groupService.CreateKeycloakAdminClient()
	if err != nil {
		return fmt.Errorf("failed to create Keycloak admin client: %w", err)
	}

	realm := s.keycloakService.realm

	// Delete the user
	err = client.DeleteUser(ctx, token, realm, userId)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
