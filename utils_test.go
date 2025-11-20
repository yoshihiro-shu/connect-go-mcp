package connectgomcp

import (
	"context"
	"regexp"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListTools(t *testing.T) {
	ctx := context.Background()

	// Create a test server with some tools
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}, nil)

	// Add test tools
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "test-tool-1",
			Description: "First test tool",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input map[string]interface{}) (*mcp.CallToolResult, interface{}, error) {
			return &mcp.CallToolResult{}, nil, nil
		},
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "test-tool-2",
			Description: "Second test tool",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input map[string]interface{}) (*mcp.CallToolResult, interface{}, error) {
			return &mcp.CallToolResult{}, nil, nil
		},
	)

	// List tools using the utility function
	tools, err := ListTools(ctx, server)
	require.NoError(t, err, "ListTools should not return an error")
	require.NotNil(t, tools, "Tools list should not be nil")

	// Verify we got the expected tools
	assert.Len(t, tools, 2, "Should have 2 tools")

	// Verify tool names and descriptions
	toolNames := make(map[string]string)
	for _, tool := range tools {
		toolNames[tool.Name] = tool.Description
	}

	assert.Equal(t, "First test tool", toolNames["test-tool-1"])
	assert.Equal(t, "Second test tool", toolNames["test-tool-2"])
}

func TestListTools_EmptyServer(t *testing.T) {
	ctx := context.Background()

	// Create a server with no tools
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "empty-server",
		Version: "1.0.0",
	}, nil)

	// List tools
	tools, err := ListTools(ctx, server)
	require.NoError(t, err, "ListTools should not return an error even for empty server")
	assert.Empty(t, tools, "Tools list should be empty")
}

func TestFilterTools(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name            string
		tools           []string
		pattern         string
		expectedRemoved []string
		expectedKept    []string
	}{
		{
			name:            "filter removes matching tools",
			tools:           []string{"GetUser", "CreateUser", "ListUsers"},
			pattern:         "^Get.*",
			expectedRemoved: []string{"GetUser"},
			expectedKept:    []string{"CreateUser", "ListUsers"},
		},
		{
			name:            "filter with List prefix",
			tools:           []string{"GetUser", "CreateUser", "ListUsers"},
			pattern:         "^List.*",
			expectedRemoved: []string{"ListUsers"},
			expectedKept:    []string{"GetUser", "CreateUser"},
		},
		{
			name:            "filter with multiple patterns",
			tools:           []string{"GetUser", "CreateUser", "ListUsers", "DeleteUser"},
			pattern:         "^(Get|Delete).*",
			expectedRemoved: []string{"GetUser", "DeleteUser"},
			expectedKept:    []string{"CreateUser", "ListUsers"},
		},
		{
			name:            "no matches - all tools kept",
			tools:           []string{"GetUser", "CreateUser", "ListUsers"},
			pattern:         "^NonExistent.*",
			expectedRemoved: []string{},
			expectedKept:    []string{"GetUser", "CreateUser", "ListUsers"},
		},
		{
			name:            "pattern matches all tools",
			tools:           []string{"GetUser", "CreateUser", "ListUsers"},
			pattern:         ".*User.*",
			expectedRemoved: []string{"GetUser", "CreateUser", "ListUsers"},
			expectedKept:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server with tools
			server := mcp.NewServer(&mcp.Implementation{
				Name:    "test-server",
				Version: "1.0.0",
			}, nil)

			// Add test tools
			for _, toolName := range tt.tools {
				mcp.AddTool(
					server,
					&mcp.Tool{
						Name:        toolName,
						Description: "Test tool: " + toolName,
					},
					func(ctx context.Context, req *mcp.CallToolRequest, input map[string]interface{}) (*mcp.CallToolResult, interface{}, error) {
						return &mcp.CallToolResult{}, nil, nil
					},
				)
			}

			// Compile pattern
			pattern, err := regexp.Compile(tt.pattern)
			require.NoError(t, err, "Pattern should compile successfully")

			// Apply filter
			FilterTools(server, pattern)

			// List tools after filtering
			tools, err := ListTools(ctx, server)
			require.NoError(t, err, "ListTools should not return an error")

			// Collect tool names
			toolNames := make(map[string]bool)
			for _, tool := range tools {
				toolNames[tool.Name] = true
			}

			// Verify removed tools are not present
			for _, removed := range tt.expectedRemoved {
				assert.False(t, toolNames[removed], "Tool %q should have been removed", removed)
			}

			// Verify kept tools are present
			for _, kept := range tt.expectedKept {
				assert.True(t, toolNames[kept], "Tool %q should have been kept", kept)
			}

			// Verify total count
			expectedCount := len(tt.expectedKept)
			assert.Equal(t, expectedCount, len(tools), "Should have %d tools after filtering", expectedCount)
		})
	}
}

func TestFilterTools_EmptyServer(t *testing.T) {
	ctx := context.Background()

	// Create an empty server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "empty-server",
		Version: "1.0.0",
	}, nil)

	// Apply filter to empty server
	pattern := regexp.MustCompile("^Get.*")
	FilterTools(server, pattern)

	// List tools - should be empty
	tools, err := ListTools(ctx, server)
	require.NoError(t, err, "ListTools should not return an error")
	assert.Empty(t, tools, "Empty server should remain empty after filtering")
}

func TestFilterTools_NilPattern(t *testing.T) {
	// Create a test server with tools
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "TestTool",
			Description: "Test tool",
		},
		func(ctx context.Context, req *mcp.CallToolRequest, input map[string]interface{}) (*mcp.CallToolResult, interface{}, error) {
			return &mcp.CallToolResult{}, nil, nil
		},
	)

	// Test that nil pattern causes panic
	assert.Panics(t, func() {
		FilterTools(server, nil)
	}, "FilterTools should panic with nil pattern")
}
