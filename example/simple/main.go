package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	connectgomcp "github.com/yoshihiro-shu/connect-go-mcp"
	userv1 "github.com/yoshihiro-shu/connect-go-mcp/example/simple/gen/user/v1"
)

func main() {
	// HTTPクライアントオプション付きのツール登録済みMCPサーバーを取得
	s := userv1.NewUserServiceMCPServer(
		"http://localhost:8080",
		connectgomcp.WithHTTPHeaders(map[string]string{
			"Authorization": "Bearer 1234567890",
		}),
	)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
