package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	connectgomcp "github.com/yoshihiro-shu/connect-go-mcp"
	userv1 "github.com/yoshihiro-shu/connect-go-mcp/example/simple/gen/user/v1"
	"github.com/yoshihiro-shu/connect-go-mcp/example/simple/gen/user/v1/userv1connect"
)

func main() {
	// HTTPクライアント設定
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 認証インターセプター（例）
	authInterceptor := connect.UnaryInterceptorFunc(
		func(next connect.UnaryFunc) connect.UnaryFunc {
			return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
				// 実際の環境では、適切な認証トークンを設定
				if token := os.Getenv("API_TOKEN"); token != "" {
					req.Header().Set("Authorization", "Bearer "+token)
				}
				return next(ctx, req)
			}
		},
	)

	// HTTPクライアントオプション付きのツール登録済みMCPサーバーを取得
	server := userv1.NewUserServiceMCPServer(
		"https://api.example.com",
		connectgomcp.WithHTTPClient(httpClient),
	)

	// サービスクライアント作成（実際のAPIサーバー通信用、デモ目的）
	_ = userv1connect.NewUserServiceClient(
		httpClient,
		"https://api.example.com", // 実際のAPI URL
		connect.WithInterceptors(authInterceptor),
	)

	// サーバー起動
	fmt.Println("MCP server created successfully!")
	fmt.Println("Available tools:")
	fmt.Println("  - GetUser ユーザー情報取得: user_id (string)")
	fmt.Println("  - CreateUser ユーザー作成: name (string), email (string)")
	fmt.Println("  - ListUsers ユーザー一覧取得: limit (number), offset (number)")
	fmt.Println("")
	fmt.Println("This is a demo - in a real application, you would:")
	fmt.Println("1. Initialize stdio transport")
	fmt.Println("2. Set up MCP transport handlers")
	fmt.Println("3. Start the MCP server loop")
	fmt.Printf("Generated server: %+v\n", server)
}

// Example of how you would use the generated client (not part of MCP server)
func exampleClientUsage() {
	// HTTPクライアント設定
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// サービスクライアント作成
	client := userv1connect.NewUserServiceClient(
		httpClient,
		"https://api.example.com",
	)

	ctx := context.Background()

	// ユーザー作成
	createResp, err := client.CreateUser(ctx, connect.NewRequest(&userv1.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}))
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("Created user: %+v\n", createResp.Msg.User)

	// ユーザー取得
	getResp, err := client.GetUser(ctx, connect.NewRequest(&userv1.GetUserRequest{
		UserId: createResp.Msg.User.Id,
	}))
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	fmt.Printf("Retrieved user: %+v\n", getResp.Msg.User)

	// ユーザー一覧取得
	listResp, err := client.ListUsers(ctx, connect.NewRequest(&userv1.ListUsersRequest{
		Limit:  10,
		Offset: 0,
	}))
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}

	fmt.Printf("Listed %d users (total: %d)\n", len(listResp.Msg.Users), listResp.Msg.Total)
}
