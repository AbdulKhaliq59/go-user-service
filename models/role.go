package models

// Role represents a Keycloak role
type Role struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ClientRole  bool   `json:"clientRole,omitempty"`
	ClientID    string `json:"clientId,omitempty"`
}
