package connectgomcp

import (
	"context"
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
