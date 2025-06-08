package main

import (
	"fmt"

	userv1mcp "github.com/yoshihiro-shu/connect-go-mcp/example/simple/gen/userv1mcp/user/v1"

	"github.com/mark3labs/mcp-go/server"
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

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
