module github.com/yoshihiro-shu/connect-go-mcp/example/simple

go 1.23.0

toolchain go1.24.1

require (
	connectrpc.com/connect v1.18.1
	github.com/modelcontextprotocol/go-sdk v0.5.0
	github.com/yoshihiro-shu/connect-go-mcp v0.0.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/google/jsonschema-go v0.2.3 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
)

// Use local connect-go-mcp
replace github.com/yoshihiro-shu/connect-go-mcp => ../..
