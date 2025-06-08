package connect_mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

type ToolHandler struct {
	config *toolConfig
}

type toolConfig struct {
	baseURL     string
	httpClient  *http.Client
	httpHeaders http.Header
}

func NewToolConfig(baseURL string) *toolConfig {
	return &toolConfig{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func NewToolHandler(baseURL string, opts ...ClientOption) *ToolHandler {
	config := NewToolConfig(baseURL)
	for _, opt := range opts {
		opt.apply(config)
	}
	return &ToolHandler{config: config}
}

func (h *ToolHandler) httpRequest(ctx context.Context, req mcp.CallToolRequest, endpoint string) (*mcp.CallToolResult, error) {
	arguments := req.Params.Arguments

	reqj, err := json.Marshal(arguments)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", h.config.baseURL, endpoint)
	hreq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqj))
	if err != nil {
		return nil, err
	}

	hreq.Header.Set("Content-Type", "application/json")

	for k, v := range h.config.httpHeaders {
		for _, vv := range v {
			hreq.Header.Add(k, vv)
		}
	}

	resp, err := h.config.httpClient.Do(hreq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result mcp.CallToolResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (h *ToolHandler) Handle(ctx context.Context, req mcp.CallToolRequest, endpoint string) (*mcp.CallToolResult, error) {
	return h.httpRequest(ctx, req, endpoint)
}
