module github.com/yoshihiro-shu/connect-go-mcp/example/simple

go 1.23

require (
	connectrpc.com/connect v1.18.1
	github.com/mark3labs/mcp-go v0.22.0
	github.com/yoshihiro-shu/connect-go-mcp v0.0.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
)

// Use local connect-go-mcp
replace github.com/yoshihiro-shu/connect-go-mcp => ../..
