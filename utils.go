package connectgomcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListTools retrieves the list of tools from an MCP server.
// This is useful for testing and debugging purposes.
//
// Example:
//
//	server := userv1mcp.NewUserServiceMCPServer("http://localhost:8080")
//	tools, err := connectgomcp.ListTools(ctx, server)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, tool := range tools {
//	    fmt.Printf("Tool: %s - %s\n", tool.Name, tool.Description)
//	}
func ListTools(ctx context.Context, server *mcp.Server) ([]*mcp.Tool, error) {
	// Create in-memory transport pair
	ct, st := mcp.NewInMemoryTransports()

	// Connect server side
	serverSession, err := server.Connect(ctx, st, nil)
	if err != nil {
		return nil, fmt.Errorf("server connect failed: %w", err)
	}
	defer serverSession.Close()

	// Create and connect client
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "connectgomcp-utils",
		Version: "1.0.0",
	}, nil)

	clientSession, err := client.Connect(ctx, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("client connect failed: %w", err)
	}
	defer clientSession.Close()

	// List tools
	result, err := clientSession.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		return nil, fmt.Errorf("list tools failed: %w", err)
	}

	return result.Tools, nil
}
