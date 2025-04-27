# connect-go-mpc

gRPCのProtocol BufferファイルからMCPServerを生成するプラグイン

## 概要

`connect-go-mpc`は、Protocol Bufferの定義ファイル（.proto）からMark3Labs MCP Serverのツール実装を自動生成するプロトコルコンパイラプラグインです。gRPCサービスの各RPCメソッドをMCPサーバーのツールとして登録し、リクエスト処理をgRPCサーバーに転送する実装を生成します。

## インストール

```bash
go install github.com/yoshihiro-shu/cmd/connect-go-mpc@latest
```

## 使用方法

```bash
protoc --connect-go-mpc_out=. \
       --connect-go-mpc_opt=module=github.com/example/myproject \
       ./proto/*.proto
```

または、bufを使用する場合:

```yaml
# buf.gen.yaml
version: v2
plugins:
  - local: protoc-gen-go
    out: gen/proto
    opt: paths=source_relative
  - local: protoc-gen-go-grpc
    out: gen/proto
    opt: paths=source_relative
  - local: connect-go-mpc
    out: gen/proto
    opt: paths=source_relative
```

そして、以下のコマンドを実行:

```bash
buf generate
```

## 機能

- Protocol Bufferからツール定義を自動生成
- RPC定義コメントをツール名に変換
- リクエスト定義コメントをツール説明に変換
- リクエストフィールドを自動的にツールパラメータに変換
- Connectクライアントとの連携機能

## コメント規約

RPC定義のコメントはツール名として使用されます：

```protobuf
// GetUser ユーザー情報取得
rpc GetUser(GetUserRequest) returns (GetUserResponse);
```

リクエスト定義のコメントはツールの説明として使用されます：

```protobuf
// GetUserRequest ユーザー情報リクエスト
// パラメータ user_id: ユーザーID（必須）
// 戻り値: ユーザー名、メールアドレス、作成日時
message GetUserRequest {
  string user_id = 1; // ユーザーID
}
```

## 生成コード例

上記のProto定義から以下のようなコードが生成されます：

```go
// NewMCPServerWithTools は設定済みの MCP サーバーを生成して返します
func NewMCPServerWithTools(client userv1connect.UserServiceClient) *mcp.Server {
  server := mcp.NewServer()
  
  server.AddTool(
    mcp.NewTool("GetUser", 
      mcp.WithDescription("ユーザー情報リクエスト パラメータ user_id: ユーザーID（必須） 戻り値: ユーザー名、メールアドレス、作成日時"),
      mcp.WithString("user_id", mcp.Required(), mcp.Description("ユーザーID")),
    ),
    func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
      return handleGetUser(ctx, req, client)
    },
  )
  
  // 他のツール登録...
  
  return server
}

func handleGetUser(ctx context.Context, req mcp.CallToolRequest, client userv1connect.UserServiceClient) (*mcp.CallToolResult, error) {
  userID, _ := req.Params.Arguments["user_id"].(string)
  
  resp, err := client.GetUser(ctx, connect.NewRequest(&userv1.GetUserRequest{
    UserId: userID,
  }))
  if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
  }
  
  return mcp.NewToolResultJSON(map[string]interface{}{
    "name": resp.Msg.User.Name,
    "email": resp.Msg.User.Email,
    "created_at": resp.Msg.User.CreatedAt.AsTime().Format(time.RFC3339),
  }), nil
}
```

## MCPサーバー起動例

```go
func main() {
  // HTTPクライアント設定
  httpClient := &http.Client{
    Timeout: 30 * time.Second,
  }
  
  // 認証インターセプター
  authInterceptor := connect.UnaryInterceptorFunc(
    func(next connect.UnaryFunc) connect.UnaryFunc {
      return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
        req.Header().Set("Authorization", "Bearer "+os.Getenv("API_TOKEN"))
        return next(ctx, req)
      }
    },
  )
  
  // サービスクライアント作成
  userClient := userv1connect.NewUserServiceClient(
    httpClient,
    "https://api.example.com",
    connect.WithInterceptors(authInterceptor),
  )
  
  // ツール登録済みのMCPサーバーを取得
  server := NewMCPServerWithTools(userClient)
  
  // サーバー起動
  log.Fatal(server.Start(":3000"))
}
```

## オプション

| オプション | 説明 | デフォルト |
|------------|------|------------|
| `module` | 生成コードのGoモジュール名 | 現在のディレクトリ |
| `out` | 出力ディレクトリ | `./mpcserver` |
| `package` | 生成コードのパッケージ名 | `mpcserver` |

## 制限事項

- 現在のバージョンではストリーミングRPCはサポートされていません
- 複雑なメッセージ型（入れ子のメッセージなど）は単純なJSONオブジェクトに変換されます

## 開発手順

### 環境構築

```bash
# リポジトリのクローン
git clone https://github.com/yoshihiro-shu/connect-go-mpc.git
cd connect-go-mpc

# 依存関係のインストール
go mod tidy

# プラグインのビルドとインストール
go install
```

### サンプルの実行

```bash
# サンプルディレクトリに移動
cd examples/simple

# サンプルコードの生成とビルド
make build

# サンプルの実行
make run
```

## 依存関係

- github.com/mark3labs/mcp-go
- golang.org/x/net
- google.golang.org/protobuf
- connectrpc.com/connect

## ライセンス

MIT

## 貢献

バグ報告、機能リクエスト、プルリクエストを歓迎します。

---

**注意**: このパッケージはプロトタイプ段階です。APIは変更される可能性があります。