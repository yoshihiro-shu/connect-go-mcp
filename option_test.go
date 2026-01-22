package connectgomcp

import (
	"net/http"
	"testing"
)

func TestWithHTTPHeaders_DoesNotOverwriteHTTPClient(t *testing.T) {
	// When only WithHTTPHeaders is used, http.DefaultClient should be preserved
	handler := NewToolHandler("http://localhost:8080", WithHTTPHeaders(map[string]string{
		"Authorization": "Bearer token",
	}))

	if handler.config.httpClient == nil {
		t.Error("httpClient should not be nil when only WithHTTPHeaders is used")
	}

	if handler.config.httpClient != http.DefaultClient {
		t.Error("httpClient should be http.DefaultClient when only WithHTTPHeaders is used")
	}

	if handler.config.httpHeaders.Get("Authorization") != "Bearer token" {
		t.Errorf("expected Authorization header to be 'Bearer token', got '%s'", handler.config.httpHeaders.Get("Authorization"))
	}
}

func TestWithHTTPClient_SetsHTTPClient(t *testing.T) {
	// When WithHTTPClient is used, the specified client should be used
	customClient := &http.Client{}
	handler := NewToolHandler("http://localhost:8080", WithHTTPClient(customClient))

	if handler.config.httpClient == nil {
		t.Error("httpClient should not be nil")
	}

	if handler.config.httpClient != customClient {
		t.Error("httpClient should be the custom client")
	}
}

func TestWithHTTPClientAndHeaders_BothAreSet(t *testing.T) {
	// When both options are used, both settings should be applied
	customClient := &http.Client{}
	handler := NewToolHandler("http://localhost:8080",
		WithHTTPClient(customClient),
		WithHTTPHeaders(map[string]string{
			"Authorization": "Bearer token",
		}),
	)

	if handler.config.httpClient == nil {
		t.Error("httpClient should not be nil")
	}

	if handler.config.httpClient != customClient {
		t.Error("httpClient should be the custom client")
	}

	if handler.config.httpHeaders.Get("Authorization") != "Bearer token" {
		t.Errorf("expected Authorization header to be 'Bearer token', got '%s'", handler.config.httpHeaders.Get("Authorization"))
	}
}

func TestWithHTTPHeadersThenClient_BothAreSet(t *testing.T) {
	// Order should not matter - test headers first, then client
	customClient := &http.Client{}
	handler := NewToolHandler("http://localhost:8080",
		WithHTTPHeaders(map[string]string{
			"X-Custom-Header": "value",
		}),
		WithHTTPClient(customClient),
	)

	if handler.config.httpClient == nil {
		t.Error("httpClient should not be nil")
	}

	if handler.config.httpClient != customClient {
		t.Error("httpClient should be the custom client")
	}

	if handler.config.httpHeaders.Get("X-Custom-Header") != "value" {
		t.Errorf("expected X-Custom-Header to be 'value', got '%s'", handler.config.httpHeaders.Get("X-Custom-Header"))
	}
}

func TestDefaultToolHandler_UsesDefaultClient(t *testing.T) {
	// With no options, http.DefaultClient should be used
	handler := NewToolHandler("http://localhost:8080")

	if handler.config.httpClient == nil {
		t.Error("httpClient should not be nil")
	}

	if handler.config.httpClient != http.DefaultClient {
		t.Error("httpClient should be http.DefaultClient by default")
	}
}
