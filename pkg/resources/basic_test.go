package resources

import (
	"strings"
	"testing"

	"github.com/granicus/ehq-mcp-server/internal/mcp"
)

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestHelloWorldResource(t *testing.T) {
	resource := HelloWorldResource()

	if resource.URI != "hello://world" {
		t.Errorf("Expected URI 'hello://world', got '%s'", resource.URI)
	}

	if resource.Name != "Hello World" {
		t.Errorf("Expected name 'Hello World', got '%s'", resource.Name)
	}

	if resource.Description == "" {
		t.Error("Expected non-empty description")
	}

	if resource.MimeType != "text/plain" {
		t.Errorf("Expected mime type 'text/plain', got '%s'", resource.MimeType)
	}
}

func TestServerInfoResource(t *testing.T) {
	resource := ServerInfoResource()

	if resource.URI != "info://server" {
		t.Errorf("Expected URI 'info://server', got '%s'", resource.URI)
	}

	if resource.Name != "Server Information" {
		t.Errorf("Expected name 'Server Information', got '%s'", resource.Name)
	}

	if resource.Description == "" {
		t.Error("Expected non-empty description")
	}

	if resource.MimeType != "text/plain" {
		t.Errorf("Expected mime type 'text/plain', got '%s'", resource.MimeType)
	}
}

func TestAllResourcesReturnValidMCPResource(t *testing.T) {
	resources := []mcp.Resource{
		HelloWorldResource(),
		ServerInfoResource(),
	}

	for _, resource := range resources {
		if resource.URI == "" {
			t.Error("Resource URI should not be empty")
		}

		if resource.Name == "" {
			t.Error("Resource name should not be empty")
		}

		if resource.MimeType == "" {
			t.Error("Resource mime type should not be empty")
		}

		// URI should have a scheme (protocol://)
		if len(resource.URI) < 3 || !contains(resource.URI, "://") {
			t.Errorf("Resource URI '%s' should have a scheme", resource.URI)
		}
	}
}

func TestResourceURIsAreUnique(t *testing.T) {
	resources := []mcp.Resource{
		HelloWorldResource(),
		ServerInfoResource(),
	}

	uris := make(map[string]bool)
	for _, resource := range resources {
		if uris[resource.URI] {
			t.Errorf("Duplicate resource URI found: %s", resource.URI)
		}
		uris[resource.URI] = true
	}
}

func TestResourceNamesAreUnique(t *testing.T) {
	resources := []mcp.Resource{
		HelloWorldResource(),
		ServerInfoResource(),
	}

	names := make(map[string]bool)
	for _, resource := range resources {
		if names[resource.Name] {
			t.Errorf("Duplicate resource name found: %s", resource.Name)
		}
		names[resource.Name] = true
	}
}
