// Package api provides EHQ API client functionality
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// EHQClient represents a client for EngagementHQ API
type EHQClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// AuthRequest represents the authentication request payload
type AuthRequest struct {
	Data AuthData `json:"data"`
}

// AuthData represents the data portion of the auth request
type AuthData struct {
	Attributes AuthAttributes `json:"attributes"`
}

// AuthAttributes represents the attributes portion of the auth request
type AuthAttributes struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Data AuthResponseData `json:"data"`
}

// AuthResponseData represents the data portion of the auth response
type AuthResponseData struct {
	Attributes AuthResponseAttributes `json:"attributes"`
}

// AuthResponseAttributes represents the attributes portion of the auth response
type AuthResponseAttributes struct {
	Token string `json:"token"`
}

// ProjectsResponse represents the JSON-API compliant projects response
type ProjectsResponse struct {
	Data []ProjectData `json:"data"`
}

// ProjectData represents a single project in JSON-API format
type ProjectData struct {
	Type       string                 `json:"type"`
	ID         string                 `json:"id,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
}

// NewEHQClient creates a new EHQ API client
func NewEHQClient(baseURL string) *EHQClient {
	// Ensure baseURL doesn't end with a slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &EHQClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Authenticate performs authentication and stores the JWT token
func (c *EHQClient) Authenticate(login, password string) error {
	authReq := AuthRequest{
		Data: AuthData{
			Attributes: AuthAttributes{
				Login:    login,
				Password: password,
			},
		},
	}

	jsonData, err := json.Marshal(authReq)
	if err != nil {
		return fmt.Errorf("failed to marshal auth request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v2/tokens", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	fmt.Printf("Auth Response: %+v\n", authResp)

	c.Token = authResp.Data.Attributes.Token
	return nil
}

// GetProjects fetches projects from the EHQ API with optional search filter
func (c *EHQClient) GetProjects(search string) (*ProjectsResponse, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("client not authenticated - call Authenticate() first")
	}

	projectURL := fmt.Sprintf("%s/api/v2/projects?filterable=true", c.BaseURL)

	// Add search filter as query parameter if provided
	if search != "" {
		projectURL += fmt.Sprintf("&filters[search]=%s", url.QueryEscape(search))
	}

	req, err := http.NewRequest("GET", projectURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create projects request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send projects request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("projects request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var projectsResp ProjectsResponse
	if err := json.NewDecoder(resp.Body).Decode(&projectsResp); err != nil {
		return nil, fmt.Errorf("failed to decode projects response: %w", err)
	}

	return &projectsResp, nil
}
