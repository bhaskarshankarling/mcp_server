package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewEHQClient(t *testing.T) {
	baseURL := "https://example.com"
	client := NewEHQClient(baseURL)

	if client.BaseURL != baseURL {
		t.Errorf("Expected base URL '%s', got '%s'", baseURL, client.BaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTP client to be initialized")
	}

	if client.Token != "" {
		t.Error("Expected token to be empty initially")
	}
}

func TestNewEHQClientTrimsSlash(t *testing.T) {
	baseURL := "https://example.com/"
	expectedURL := "https://example.com"
	client := NewEHQClient(baseURL)

	if client.BaseURL != expectedURL {
		t.Errorf("Expected base URL '%s', got '%s'", expectedURL, client.BaseURL)
	}
}

func TestAuthenticate(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v2/tokens" {
			t.Errorf("Expected path '/api/v2/tokens', got '%s'", r.URL.Path)
		}

		// Check content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		// Parse request body
		var authReq AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
			t.Fatalf("Failed to decode auth request: %v", err)
		}

		// Check request data
		if authReq.Data.Attributes.Login != "testuser" {
			t.Errorf("Expected login 'testuser', got '%s'", authReq.Data.Attributes.Login)
		}

		if authReq.Data.Attributes.Password != "testpass" {
			t.Errorf("Expected password 'testpass', got '%s'", authReq.Data.Attributes.Password)
		}

		// Send response
		response := AuthResponse{
			Data: AuthResponseData{
				Attributes: AuthResponseAttributes{
					Token: "test-token-123",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewEHQClient(server.URL)
	err := client.Authenticate("testuser", "testpass")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client.Token != "test-token-123" {
		t.Errorf("Expected token 'test-token-123', got '%s'", client.Token)
	}
}

func TestAuthenticateFailure(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}))
	defer server.Close()

	client := NewEHQClient(server.URL)
	err := client.Authenticate("baduser", "badpass")

	if err == nil {
		t.Error("Expected error for failed authentication")
	}

	if client.Token != "" {
		t.Error("Expected token to remain empty on failed authentication")
	}
}

func TestGetProjects(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		if r.URL.Path != "/api/v2/projects" {
			t.Errorf("Expected path '/api/v2/projects', got '%s'", r.URL.Path)
		}

		// Check authorization header
		expectedAuth := "Bearer test-token"
		if r.Header.Get("Authorization") != expectedAuth {
			t.Errorf("Expected Authorization '%s', got '%s'", expectedAuth, r.Header.Get("Authorization"))
		}

		// Check query parameters
		if !r.URL.Query().Has("filterable") {
			t.Error("Expected 'filterable' query parameter")
		}

		// Send response
		response := ProjectsResponse{
			Data: []ProjectData{
				{
					Type: "projects",
					ID:   "1",
					Attributes: map[string]interface{}{
						"name":        "Test Project 1",
						"description": "A test project",
					},
				},
				{
					Type: "projects",
					ID:   "2",
					Attributes: map[string]interface{}{
						"name":        "Test Project 2",
						"description": "Another test project",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewEHQClient(server.URL)
	client.Token = "test-token" // Set token directly for testing

	projects, err := client.GetProjects("")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(projects.Data) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects.Data))
	}

	if projects.Data[0].ID != "1" {
		t.Errorf("Expected first project ID '1', got '%s'", projects.Data[0].ID)
	}

	if projects.Data[0].Type != "projects" {
		t.Errorf("Expected first project type 'projects', got '%s'", projects.Data[0].Type)
	}

	// Check attributes
	name, ok := projects.Data[0].Attributes["name"].(string)
	if !ok || name != "Test Project 1" {
		t.Errorf("Expected first project name 'Test Project 1', got '%s'", name)
	}
}

func TestGetProjectsWithSearch(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check query parameters
		if !r.URL.Query().Has("filterable") {
			t.Error("Expected 'filterable' query parameter")
		}

		searchParam := r.URL.Query().Get("filters[search]")
		if searchParam != "test search" {
			t.Errorf("Expected search parameter 'test search', got '%s'", searchParam)
		}

		// Send empty response
		response := ProjectsResponse{Data: []ProjectData{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewEHQClient(server.URL)
	client.Token = "test-token"

	_, err := client.GetProjects("test search")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGetProjectsNotAuthenticated(t *testing.T) {
	client := NewEHQClient("https://example.com")
	// Don't set token

	_, err := client.GetProjects("")
	if err == nil {
		t.Error("Expected error when not authenticated")
	}

	expectedErrorMsg := "client not authenticated - call Authenticate() first"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}

func TestGetProjectsFailure(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewEHQClient(server.URL)
	client.Token = "test-token"

	_, err := client.GetProjects("")
	if err == nil {
		t.Error("Expected error for failed request")
	}
}

func TestAuthRequestStructure(t *testing.T) {
	authReq := AuthRequest{
		Data: AuthData{
			Attributes: AuthAttributes{
				Login:    "testuser",
				Password: "testpass",
			},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(authReq)
	if err != nil {
		t.Fatalf("Failed to marshal auth request: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled AuthRequest
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal auth request: %v", err)
	}

	if unmarshaled.Data.Attributes.Login != "testuser" {
		t.Errorf("Expected login 'testuser', got '%s'", unmarshaled.Data.Attributes.Login)
	}

	if unmarshaled.Data.Attributes.Password != "testpass" {
		t.Errorf("Expected password 'testpass', got '%s'", unmarshaled.Data.Attributes.Password)
	}
}

func TestProjectsResponseStructure(t *testing.T) {
	projectsResp := ProjectsResponse{
		Data: []ProjectData{
			{
				Type: "projects",
				ID:   "1",
				Attributes: map[string]interface{}{
					"name":        "Test Project",
					"description": "A test project",
					"status":      "active",
				},
			},
		},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(projectsResp)
	if err != nil {
		t.Fatalf("Failed to marshal projects response: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled ProjectsResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal projects response: %v", err)
	}

	if len(unmarshaled.Data) != 1 {
		t.Errorf("Expected 1 project, got %d", len(unmarshaled.Data))
	}

	project := unmarshaled.Data[0]
	if project.Type != "projects" {
		t.Errorf("Expected type 'projects', got '%s'", project.Type)
	}

	if project.ID != "1" {
		t.Errorf("Expected ID '1', got '%s'", project.ID)
	}

	name, ok := project.Attributes["name"].(string)
	if !ok || name != "Test Project" {
		t.Errorf("Expected name 'Test Project', got '%s'", name)
	}
}
