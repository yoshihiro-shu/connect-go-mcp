package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	connectgomcp "github.com/yoshihiro-shu/connect-go-mcp"
	"github.com/yoshihiro-shu/connect-go-mcp/cmd/protoc-gen-connect-go-mcp/testdata/greet/gen/greetv1mcp"
)

// setupMockBackend creates a mock HTTP server for testing
func setupMockBackend(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

// setupMCPConnection creates an MCP client-server connection for testing
func setupMCPConnection(t *testing.T, mcpServer *mcp.Server) (*mcp.ClientSession, *mcp.ServerSession) {
	t.Helper()
	ctx := context.Background()

	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	serverSession, err := mcpServer.Connect(ctx, serverTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { serverSession.Close() })

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	clientSession, err := client.Connect(ctx, clientTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { clientSession.Close() })

	return clientSession, serverSession
}

// TestIntegrationMCPServerCreation verifies that the generated MCP server can be instantiated
func TestIntegrationMCPServerCreation(t *testing.T) {
	mockServer := setupMockBackend(t, func(w http.ResponseWriter, r *http.Request) {})
	mcpServer := greetv1mcp.NewGreetServiceMCPServer(mockServer.URL)
	assert.NotNil(t, mcpServer, "MCP server should not be nil")
}

// TestIntegrationMCPServerToolListing verifies that tools are properly registered and can be listed
func TestIntegrationMCPServerToolListing(t *testing.T) {
	mockServer := setupMockBackend(t, func(w http.ResponseWriter, r *http.Request) {})
	mcpServer := greetv1mcp.NewGreetServiceMCPServer(mockServer.URL)
	clientSession, _ := setupMCPConnection(t, mcpServer)

	ctx := context.Background()
	result, err := clientSession.ListTools(ctx, nil)
	require.NoError(t, err)

	// Verify 2 tools are registered (Greet RPC and Ping RPC)
	assert.Len(t, result.Tools, 2, "Expected 2 tools to be registered")

	// Collect tool names
	toolNames := make(map[string]bool)
	for _, tool := range result.Tools {
		toolNames[tool.Name] = true
	}

	// Verify expected tools are present
	assert.True(t, toolNames["Greet RPC"], "Greet RPC tool should be registered")
	assert.True(t, toolNames["Ping RPC"], "Ping RPC tool should be registered")

	// Verify tool descriptions and schemas
	for _, tool := range result.Tools {
		if tool.Name == "Greet RPC" {
			assert.Contains(t, tool.Description, "Greeting request")
			assert.NotNil(t, tool.InputSchema)
			assert.Equal(t, "object", tool.InputSchema.Type)
			assert.Contains(t, tool.InputSchema.Properties, "name")
		}
		if tool.Name == "Ping RPC" {
			assert.Contains(t, tool.Description, "Ping request")
			assert.NotNil(t, tool.InputSchema)
			assert.Equal(t, "object", tool.InputSchema.Type)
			assert.Contains(t, tool.InputSchema.Properties, "message")
		}
	}
}

// TestIntegrationMCPServerToolInvocation verifies that tool calls make proper HTTP requests
func TestIntegrationMCPServerToolInvocation(t *testing.T) {
	var mu sync.Mutex
	var receivedRequest struct {
		Method      string
		Path        string
		ContentType string
		Body        map[string]any
	}

	mockServer := setupMockBackend(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		receivedRequest.Method = r.Method
		receivedRequest.Path = r.URL.Path
		receivedRequest.ContentType = r.Header.Get("Content-Type")
		json.NewDecoder(r.Body).Decode(&receivedRequest.Body)

		response := &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Hello, World!"}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	mcpServer := greetv1mcp.NewGreetServiceMCPServer(mockServer.URL)
	clientSession, _ := setupMCPConnection(t, mcpServer)

	ctx := context.Background()
	_, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "Greet RPC",
		Arguments: map[string]any{"name": "World"},
	})
	require.NoError(t, err)

	// Verify HTTP request details
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, http.MethodPost, receivedRequest.Method, "HTTP method should be POST")
	assert.Equal(t, "/Greet", receivedRequest.Path, "HTTP path should be /Greet")
	assert.Equal(t, "application/json", receivedRequest.ContentType, "Content-Type should be application/json")
	assert.Equal(t, "World", receivedRequest.Body["name"], "Request body should contain the name parameter")
}

// TestIntegrationMCPServerToolResponse verifies that responses are correctly returned
func TestIntegrationMCPServerToolResponse(t *testing.T) {
	expectedText := "Hello, TestUser!"

	mockServer := setupMockBackend(t, func(w http.ResponseWriter, r *http.Request) {
		response := &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: expectedText}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	mcpServer := greetv1mcp.NewGreetServiceMCPServer(mockServer.URL)
	clientSession, _ := setupMCPConnection(t, mcpServer)

	ctx := context.Background()
	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "Greet RPC",
		Arguments: map[string]any{"name": "TestUser"},
	})
	require.NoError(t, err)

	// Verify response
	assert.False(t, result.IsError, "Result should not be an error")
	require.NotEmpty(t, result.Content, "Result content should not be empty")

	// Verify text content
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok, "Content should be TextContent")
	assert.Equal(t, expectedText, textContent.Text, "Response text should match expected")
}

// TestIntegrationMCPServerWithCustomHeaders verifies that custom HTTP headers are sent
func TestIntegrationMCPServerWithCustomHeaders(t *testing.T) {
	var mu sync.Mutex
	var receivedHeaders http.Header

	mockServer := setupMockBackend(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		receivedHeaders = r.Header.Clone()
		mu.Unlock()

		response := &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "OK"}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	mcpServer := greetv1mcp.NewGreetServiceMCPServer(
		mockServer.URL,
		connectgomcp.WithHTTPHeaders(map[string]string{
			"Authorization":   "Bearer test-token",
			"X-Custom-Header": "custom-value",
		}),
	)
	clientSession, _ := setupMCPConnection(t, mcpServer)

	ctx := context.Background()
	_, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "Greet RPC",
		Arguments: map[string]any{"name": "World"},
	})
	require.NoError(t, err)

	// Verify custom headers were sent
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, "Bearer test-token", receivedHeaders.Get("Authorization"), "Authorization header should be set")
	assert.Equal(t, "custom-value", receivedHeaders.Get("X-Custom-Header"), "Custom header should be set")
}

// TestIntegrationMCPServerBackendError verifies error handling when backend returns an error
func TestIntegrationMCPServerBackendError(t *testing.T) {
	mockServer := setupMockBackend(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	mcpServer := greetv1mcp.NewGreetServiceMCPServer(mockServer.URL)
	clientSession, _ := setupMCPConnection(t, mcpServer)

	ctx := context.Background()
	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "Greet RPC",
		Arguments: map[string]any{"name": "World"},
	})

	// The ToolHandler returns errors either as protocol error or embedded in result
	// Based on tool_handler.go, non-200 status returns an error which becomes IsError=true
	if err == nil {
		assert.True(t, result.IsError, "Result should indicate an error")
	}
}

// TestIntegrationMCPServerPingTool verifies the Ping RPC tool works correctly
func TestIntegrationMCPServerPingTool(t *testing.T) {
	var mu sync.Mutex
	var receivedPath string
	var receivedBody map[string]any

	mockServer := setupMockBackend(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		receivedPath = r.URL.Path
		json.NewDecoder(r.Body).Decode(&receivedBody)
		mu.Unlock()

		response := &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "pong"}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	mcpServer := greetv1mcp.NewGreetServiceMCPServer(mockServer.URL)
	clientSession, _ := setupMCPConnection(t, mcpServer)

	ctx := context.Background()
	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "Ping RPC",
		Arguments: map[string]any{"message": "hello"},
	})
	require.NoError(t, err)

	// Verify correct endpoint was called
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, "/Ping", receivedPath, "HTTP path should be /Ping")
	assert.Equal(t, "hello", receivedBody["message"], "Request body should contain the message parameter")

	// Verify response
	assert.False(t, result.IsError)
	require.NotEmpty(t, result.Content)
	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok)
	assert.Equal(t, "pong", textContent.Text)
}
