package main

import (
	"context"
	"fmt"

	userv1mcp "github.com/yoshihiro-shu/connect-go-mcp/example/simple/gen/userv1mcp/user/v1"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	connectgomcp "github.com/yoshihiro-shu/connect-go-mcp"
)

func main() {
	// HTTPクライアントオプション付きのツール登録済みMCPサーバーを取得
	s := userv1mcp.NewUserServiceMCPServer(
		"http://localhost:8080",
		connectgomcp.WithHTTPHeaders(map[string]string{
			"Authorization": "Bearer 1234567890",
		}),
	)

	ctx := context.Background()
	if err := s.Run(ctx, &mcp.StdioTransport{}); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
