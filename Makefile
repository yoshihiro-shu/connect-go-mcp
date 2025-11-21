.PHONY: help test build lint clean install generate coverage ci-test

# デフォルトターゲット
help:
	@echo "Available targets:"
	@echo "  make test          - Run all tests"
	@echo "  make test-race     - Run tests with race detector"
	@echo "  make build         - Build the plugin"
	@echo "  make install       - Install the plugin locally"
	@echo "  make generate      - Generate test data"
	@echo "  make lint          - Run linters"
	@echo "  make coverage      - Generate coverage report"
	@echo "  make ci-test       - Run CI tests locally (same as GitHub Actions)"
	@echo "  make clean         - Clean build artifacts"

# テスト実行
test:
	@echo "Running tests..."
	go test -v ./...

# race detectorを有効にしてテスト
test-race:
	@echo "Running tests with race detector..."
	go test -v -race ./...

# プラグインのビルド
build:
	@echo "Building plugin..."
	go build -v -o protoc-gen-connect-go-mcp ./cmd/protoc-gen-connect-go-mcp
	@echo "Build successful: ./protoc-gen-connect-go-mcp"

# プラグインのローカルインストール
install:
	@echo "Installing plugin..."
	go build -o ~/.local/bin/protoc-gen-connect-go-mcp ./cmd/protoc-gen-connect-go-mcp
	@echo "Installed to: ~/.local/bin/protoc-gen-connect-go-mcp"

# テストデータの生成
generate:
	@echo "Generating test data..."
	cd cmd/protoc-gen-connect-go-mcp && ./generate.sh

# 依存関係の検証
verify:
	@echo "Verifying dependencies..."
	go mod verify

# Lintの実行
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not installed. Install with: brew install golangci-lint"; \
		exit 1; \
	fi

# カバレッジレポートの生成
coverage:
	@echo "Generating coverage report..."
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"
	go tool cover -func=coverage.txt | tail -1

# CI テスト（GitHub Actions と同じ）
ci-test: verify build generate test-race
	@echo ""
	@echo "✅ All CI checks passed!"

# クリーンアップ
clean:
	@echo "Cleaning build artifacts..."
	rm -f protoc-gen-connect-go-mcp
	rm -f coverage.txt coverage.html
	@echo "Clean complete"

# すべてのチェック（プッシュ前に実行推奨）
pre-push: clean ci-test lint
	@echo ""
	@echo "🎉 All checks passed! Ready to push."
