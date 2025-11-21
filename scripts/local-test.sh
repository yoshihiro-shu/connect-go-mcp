#!/bin/bash

# ローカルでGitHub Actionsのテストステップを再現するスクリプト

set -e  # エラーで停止

echo "==> Go version"
go version

echo ""
echo "==> Verify dependencies"
go mod verify

echo ""
echo "==> Build plugin"
go build -v ./cmd/protoc-gen-connect-go-mcp
./protoc-gen-connect-go-mcp --version

echo ""
echo "==> Generate test data"
cd cmd/protoc-gen-connect-go-mcp
./generate.sh
cd ../..

echo ""
echo "==> Run tests with race detector and coverage"
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

echo ""
echo "==> Coverage report"
go tool cover -func=coverage.txt | tail -1

echo ""
echo "✅ All checks passed!"
