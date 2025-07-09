// Package mcp provides the core MCP server implementation
package mcp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/granicus/ehq-mcp-server/pkg/api"
	"github.com/sirupsen/logrus"
)

// Server represents an MCP server
type Server struct {
	serverInfo   ServerInfo
	capabilities ServerCapabilities
	tools        map[string]Tool
	resources    map[string]Resource
	logger       *logrus.Logger
	handlers     map[string]func(*Message) (*Message, error)
}

// NewServer creates a new MCP server
func NewServer(name, version string) *Server {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	server := &Server{
		serverInfo: ServerInfo{
			Name:    name,
			Version: version,
		},
		capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
			Resources: &ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
			Logging: &LoggingCapability{},
		},
		tools:     make(map[string]Tool),
		resources: make(map[string]Resource),
		logger:    logger,
		handlers:  make(map[string]func(*Message) (*Message, error)),
	}

	// Register default handlers
	server.registerHandlers()
	return server
}

// registerHandlers registers the default MCP method handlers
func (s *Server) registerHandlers() {
	s.handlers["initialize"] = s.handleInitialize
	s.handlers["tools/list"] = s.handleToolsList
	s.handlers["tools/call"] = s.handleToolsCall
	s.handlers["resources/list"] = s.handleResourcesList
	s.handlers["resources/read"] = s.handleResourcesRead
}

// RegisterTool registers a new tool with the server
func (s *Server) RegisterTool(tool Tool) {
	s.tools[tool.Name] = tool
	s.logger.Infof("Registered tool: %s", tool.Name)
}

// RegisterResource registers a new resource with the server
func (s *Server) RegisterResource(resource Resource) {
	s.resources[resource.URI] = resource
	s.logger.Infof("Registered resource: %s", resource.URI)
}

// HandleMessage processes an incoming MCP message and returns a response
func (s *Server) HandleMessage(data []byte) ([]byte, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		s.logger.Errorf("Failed to parse message: %v", err)
		errorMsg := NewError(nil, ParseError, "Parse error", nil)
		return json.Marshal(errorMsg)
	}

	s.logger.Debugf("Received message: %s", msg.Method)

	// Handle the message
	response, err := s.processMessage(&msg)
	if err != nil {
		s.logger.Errorf("Error processing message: %v", err)
		errorMsg := NewError(msg.ID, InternalError, err.Error(), nil)
		return json.Marshal(errorMsg)
	}

	return json.Marshal(response)
}

// processMessage routes the message to the appropriate handler
func (s *Server) processMessage(msg *Message) (*Message, error) {
	handler, exists := s.handlers[msg.Method]
	if !exists {
		return NewError(msg.ID, MethodNotFound, fmt.Sprintf("Method not found: %s", msg.Method), nil), nil
	}

	return handler(msg)
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(msg *Message) (*Message, error) {
	s.logger.Info("Handling initialize request")

	response := InitializeResponse{
		ProtocolVersion: MCPVersion,
		Capabilities:    s.capabilities,
		ServerInfo:      s.serverInfo,
	}

	return NewResponse(msg.ID, response), nil
}

// handleToolsList handles the tools/list request
func (s *Server) handleToolsList(msg *Message) (*Message, error) {
	s.logger.Info("Handling tools/list request")

	tools := make([]Tool, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool)
	}

	response := ToolsListResponse{
		Tools: tools,
	}

	return NewResponse(msg.ID, response), nil
}

// handleToolsCall handles the tools/call request
func (s *Server) handleToolsCall(msg *Message) (*Message, error) {
	s.logger.Info("Handling tools/call request")

	// Parse the request
	var req ToolsCallRequest
	if paramsBytes, err := json.Marshal(msg.Params); err != nil {
		return NewError(msg.ID, InvalidParams, "Invalid parameters", nil), nil
	} else if err := json.Unmarshal(paramsBytes, &req); err != nil {
		return NewError(msg.ID, InvalidParams, "Invalid parameters", nil), nil
	}

	// Check if tool exists
	tool, exists := s.tools[req.Name]
	if !exists {
		return NewError(msg.ID, MethodNotFound, fmt.Sprintf("Tool not found: %s", req.Name), nil), nil
	}

	s.logger.Infof("Executing tool: %s", tool.Name)

	// Execute the tool (this is where you'd implement actual tool logic)
	result := s.executeTool(req.Name, req.Arguments)

	response := ToolsCallResponse{
		Content: []Content{
			{
				Type: "text",
				Text: result,
			},
		},
	}

	return NewResponse(msg.ID, response), nil
}

// handleResourcesList handles the resources/list request
func (s *Server) handleResourcesList(msg *Message) (*Message, error) {
	s.logger.Info("Handling resources/list request")

	resources := make([]Resource, 0, len(s.resources))
	for _, resource := range s.resources {
		resources = append(resources, resource)
	}

	response := ResourcesListResponse{
		Resources: resources,
	}

	return NewResponse(msg.ID, response), nil
}

// handleResourcesRead handles the resources/read request
func (s *Server) handleResourcesRead(msg *Message) (*Message, error) {
	s.logger.Info("Handling resources/read request")

	// Parse the request
	var req ResourcesReadRequest
	if paramsBytes, err := json.Marshal(msg.Params); err != nil {
		return NewError(msg.ID, InvalidParams, "Invalid parameters", nil), nil
	} else if err := json.Unmarshal(paramsBytes, &req); err != nil {
		return NewError(msg.ID, InvalidParams, "Invalid parameters", nil), nil
	}

	// Check if resource exists
	resource, exists := s.resources[req.URI]
	if !exists {
		return NewError(msg.ID, MethodNotFound, fmt.Sprintf("Resource not found: %s", req.URI), nil), nil
	}

	s.logger.Infof("Reading resource: %s", resource.URI)

	// Read the resource (this is where you'd implement actual resource reading)
	content := s.readResource(req.URI)

	response := ResourcesReadResponse{
		Contents: []ResourceContent{
			{
				URI:      req.URI,
				MimeType: resource.MimeType,
				Text:     content,
			},
		},
	}

	return NewResponse(msg.ID, response), nil
}

// executeTool executes a tool with the given arguments
func (s *Server) executeTool(name string, args map[string]interface{}) string {
	switch name {
	case "hello_world":
		return s.ExecuteHelloWorld(args)
	case "echo":
		return s.ExecuteEcho(args)
	case "get_time":
		return s.ExecuteGetTime(args)
	case "get_projects":
		response, _ := s.ExecuteGetProjects(args)
		jsonResponse, _ := json.Marshal(response)
		return string(jsonResponse)
	default:
		return fmt.Sprintf("Tool '%s' not implemented", name)
	}
}

// ExecuteHelloWorld executes the hello world tool
func (s *Server) ExecuteHelloWorld(args map[string]interface{}) string {
	name := "World"
	if nameArg, ok := args["name"].(string); ok && nameArg != "" {
		name = nameArg
	}

	return fmt.Sprintf("Hello, %s! 🌍\nWelcome to the EHQ MCP Server!\nCurrent time: %s",
		name, time.Now().Format("2006-01-02 15:04:05"))
}

// ExecuteEcho executes the echo tool
func (s *Server) ExecuteEcho(args map[string]interface{}) string {
	message, ok := args["message"].(string)
	if !ok {
		return "Error: message parameter is required"
	}

	return fmt.Sprintf("Echo: %s", message)
}

// ExecuteGetTime executes the get time tool
func (s *Server) ExecuteGetTime(args map[string]interface{}) string {
	format := "RFC3339"
	if formatArg, ok := args["format"].(string); ok && formatArg != "" {
		format = formatArg
	}

	now := time.Now()

	switch format {
	case "Unix":
		return fmt.Sprintf("Unix timestamp: %d", now.Unix())
	case "Kitchen":
		return fmt.Sprintf("Time: %s", now.Format(time.Kitchen))
	case "RFC3339":
		fallthrough
	default:
		return fmt.Sprintf("Time: %s", now.Format(time.RFC3339))
	}
}

// ExecuteGetProjects fetches projects from the EHQ API
func (s *Server) ExecuteGetProjects(args map[string]interface{}) (map[string]interface{}, error) {
	// Extract search parameter if provided
	search := ""
	if searchArg, ok := args["search"].(string); ok {
		search = searchArg
	}

	// Create a new EHQ API client
	client := api.NewEHQClient("https://dev.ehq.test")

	// Authenticate with the API (hardcoded credentials for now)
	err := client.Authenticate("btt_admin", "Kmcdka09")
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Authentication failed: %v", err),
		}, nil
	}

	// Fetch projects with optional search filter
	projects, err := client.GetProjects(search)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Failed to fetch projects: %v", err),
		}, nil
	}

	// Log projects for debugging
	s.logger.Infof("Projects fetched successfully: %+v", projects)

	return map[string]interface{}{
		"success": true,
		"data":    projects.Data,
		"count":   len(projects.Data),
	}, nil
}

// readResource reads a resource by URI
func (s *Server) readResource(uri string) string {
	switch uri {
	case "hello://world":
		return "Hello, World! This is a sample resource from the MCP server."
	case "info://server":
		return fmt.Sprintf("Server: %s v%s\nProtocol: %s",
			s.serverInfo.Name, s.serverInfo.Version, MCPVersion)
	default:
		return fmt.Sprintf("Resource '%s' content not available", uri)
	}
}

// SetLogLevel sets the logging level
func (s *Server) SetLogLevel(level logrus.Level) {
	s.logger.SetLevel(level)
}
