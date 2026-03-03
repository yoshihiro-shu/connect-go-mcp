package connectgomcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestHttpRequest_NonOK_WithConnectError(t *testing.T) {
	wireErr := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    "not_found",
		Message: "the requested resource was not found",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(wireErr)
	}))
	defer server.Close()

	handler := NewToolHandler(server.URL)
	result, err := handler.Handle(context.Background(), &mcp.CallToolRequest{}, "test.Service/Method", map[string]any{})

	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if !result.IsError {
		t.Error("expected IsError to be true")
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Content))
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("expected TextContent")
	}
	if tc.Text != wireErr.Message {
		t.Errorf("expected message %q, got %q", wireErr.Message, tc.Text)
	}
}

func TestHttpRequest_NonOK_WithoutConnectError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	handler := NewToolHandler(server.URL)
	result, err := handler.Handle(context.Background(), &mcp.CallToolRequest{}, "test.Service/Method", map[string]any{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if result != nil {
		t.Error("expected nil result")
	}
}

func TestHttpRequest_NonOK_EmptyMessage(t *testing.T) {
	wireErr := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    "internal",
		Message: "",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(wireErr)
	}))
	defer server.Close()

	handler := NewToolHandler(server.URL)
	result, err := handler.Handle(context.Background(), &mcp.CallToolRequest{}, "test.Service/Method", map[string]any{})

	if err == nil {
		t.Fatal("expected error for empty message, got nil")
	}
	if result != nil {
		t.Error("expected nil result")
	}
}
