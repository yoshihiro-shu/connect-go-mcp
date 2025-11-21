#!/bin/bash

go build ./cmd/protoc-gen-connect-go-mcp
mv ./protoc-gen-connect-go-mcp ~/.local/bin/protoc-gen-connect-go-mcp
