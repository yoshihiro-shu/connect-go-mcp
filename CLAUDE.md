# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `protoc-gen-connect-go-mcp`, a Protocol Buffers compiler plugin that generates MCP (Model Context Protocol) server implementations from gRPC service definitions. It bridges gRPC services with the MCP ecosystem, enabling tools like Claude Desktop to interact with gRPC services.

## Common Development Commands

### Building and Installing the Plugin
```bash
# Build the plugin
go build ./cmd/protoc-gen-connect-go-mcp

# Install locally (for development)
./scripts/deploy_local.sh

# Install globally
go install github.com/yoshihiro-shu/connect-go-mcp/cmd/protoc-gen-connect-go-mcp@latest
```

### Running Tests
```bash
# Run all tests
cd cmd/protoc-gen-connect-go-mcp && go test -v ./...

# Run specific test
go test -v -run TestGenerate ./cmd/protoc-gen-connect-go-mcp

# Regenerate test data (after modifying test protos)
cd cmd/protoc-gen-connect-go-mcp && ./generate.sh
```

### Working with the Example
```bash
cd example/simple

# Generate code from proto files
make generate

# Build the example
make build

# Run the example server
make run

# Clean generated files
make clean

# Check if everything compiles
make check
```

## High-Level Architecture

### Core Components

1. **Protoc Plugin (`cmd/protoc-gen-connect-go-mcp/`)**
   - `main.go`: Entry point that implements the protoc plugin interface
   - `generator/`: Core code generation logic that produces MCP server implementations
   - `parser/`: Proto file parsing and analysis
   - `template/`: Go templates for code generation

2. **Library Components (root directory)**
   - `tool_handler.go`: Base implementation for MCP tool handlers
   - `option.go`: Client configuration options for generated servers

3. **Generated Code Pattern**
   - Input: `.proto` files defining gRPC services
   - Output: `*mcp/*.mcpserver.go` files containing MCP server implementations
   - Each RPC method becomes an MCP tool with parameter mappings

### Code Generation Flow

1. **Proto Parsing**: The plugin receives compiled proto descriptors from protoc
2. **Service Analysis**: Extracts services, methods, and message types
3. **Tool Generation**: Creates MCP tool definitions for each RPC method
4. **Server Generation**: Produces a complete MCP server implementation that:
   - Converts MCP tool calls to gRPC requests
   - Handles parameter validation and type conversion
   - Forwards requests to the actual gRPC server via HTTP

### Key Design Decisions

- **Package Organization**: Generated code uses configurable package suffix (default: `mcp`)
- **Comment Preservation**: Proto comments become MCP tool descriptions
- **Type Safety**: Maintains strong typing through generated parameter structs
- **HTTP Transport**: Uses Connect-Go's HTTP transport for gRPC communication

## Testing Strategy

- **Integration Tests**: Verify end-to-end code generation matches expected outputs
- **Unit Tests**: Test individual components (generator, parser)
- **Test Data**: Located in `testdata/` with pre-generated expected outputs
- **CI/CD**: GitHub Actions runs tests on PRs and manual triggers

## Important Notes

- When modifying the code generator, always run tests to ensure compatibility
- Generated code depends on both Connect-Go and the official MCP Go SDK (github.com/modelcontextprotocol/go-sdk v0.5.0)
- The plugin supports both buf and protoc workflows
- Package suffix parameter allows flexible code organization
- Multi-line comments in proto files are preserved with line breaks (\n) in tool descriptions
- Package name has been standardized to `connectgomcp` following Go naming conventions
- When updating dependencies, ensure test fixtures in testdata/ are updated accordingly
