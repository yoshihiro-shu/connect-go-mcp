package connect_mcpserver

import (
	"net/http"
)

type ClientOption interface {
	apply(*toolConfig)
}

type httpClientOption struct {
	client  *http.Client
	headers map[string]string
}

func (o *httpClientOption) apply(c *toolConfig) {
	c.httpClient = o.client

	if o.headers != nil {
		c.httpHeaders = http.Header{}
		for k, v := range o.headers {
			c.httpHeaders.Set(k, v)
		}
	}
}

func WithHTTPClient(client *http.Client) ClientOption {
	return &httpClientOption{client: client}
}

func WithHTTPHeaders(headers map[string]string) ClientOption {
	return &httpClientOption{headers: headers}
}
