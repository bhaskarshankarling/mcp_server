// Package mcp provides the core Model Context Protocol implementation
package mcp

import (
	"encoding/json"
	"fmt"
)

// MCPVersion represents the MCP protocol version
const MCPVersion = "2024-11-05"

// Message represents a base MCP message
type Message struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents an MCP error
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// InitializeRequest represents the initialize request
type InitializeRequest struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    Capabilities `json:"capabilities"`
	ClientInfo      ClientInfo   `json:"clientInfo"`
}

// InitializeResponse represents the initialize response
type InitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

// Capabilities represents client capabilities
type Capabilities struct {
	Roots    *RootsCapability    `json:"roots,omitempty"`
	Sampling *SamplingCapability `json:"sampling,omitempty"`
}

// ServerCapabilities represents server capabilities
type ServerCapabilities struct {
	Logging   *LoggingCapability   `json:"logging,omitempty"`
	Prompts   *PromptsCapability   `json:"prompts,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Tools     *ToolsCapability     `json:"tools,omitempty"`
}

// RootsCapability represents roots capability
type RootsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// SamplingCapability represents sampling capability
type SamplingCapability struct{}

// LoggingCapability represents logging capability
type LoggingCapability struct{}

// PromptsCapability represents prompts capability
type PromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ResourcesCapability represents resources capability
type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

// ToolsCapability represents tools capability
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ClientInfo represents client information
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ServerInfo represents server information
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// ToolsListRequest represents a tools/list request
type ToolsListRequest struct{}

// ToolsListResponse represents a tools/list response
type ToolsListResponse struct {
	Tools []Tool `json:"tools"`
}

// ToolsCallRequest represents a tools/call request
type ToolsCallRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ToolsCallResponse represents a tools/call response
type ToolsCallResponse struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content represents content in responses
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Resource represents an MCP resource
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// ResourcesListRequest represents a resources/list request
type ResourcesListRequest struct{}

// ResourcesListResponse represents a resources/list response
type ResourcesListResponse struct {
	Resources []Resource `json:"resources"`
}

// ResourcesReadRequest represents a resources/read request
type ResourcesReadRequest struct {
	URI string `json:"uri"`
}

// ResourcesReadResponse represents a resources/read response
type ResourcesReadResponse struct {
	Contents []ResourceContent `json:"contents"`
}

// ResourceContent represents resource content
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text,omitempty"`
	Blob     string `json:"blob,omitempty"`
}

// NewMessage creates a new MCP message
func NewMessage(method string, params interface{}) *Message {
	return &Message{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
}

// NewResponse creates a new MCP response
func NewResponse(id interface{}, result interface{}) *Message {
	return &Message{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// NewError creates a new MCP error response
func NewError(id interface{}, code int, message string, data interface{}) *Message {
	return &Message{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// String returns a string representation of the message
func (m *Message) String() string {
	data, _ := json.Marshal(m)
	return string(data)
}

// Error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// ErrorMessages maps error codes to messages
var ErrorMessages = map[int]string{
	ParseError:     "Parse error",
	InvalidRequest: "Invalid Request",
	MethodNotFound: "Method not found",
	InvalidParams:  "Invalid params",
	InternalError:  "Internal error",
}

// GetErrorMessage returns the error message for a given code
func GetErrorMessage(code int) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return fmt.Sprintf("Unknown error (%d)", code)
}
