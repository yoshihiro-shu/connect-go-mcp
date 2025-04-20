package connect_mpcserver

import "net/http"

type ClientOption interface {
	apply(*toolConfig)
}

type httpClientOption struct {
	client *http.Client
}

func (o *httpClientOption) apply(c *toolConfig) {
	c.httpClient = o.client
}

func WithHTTPClient(client *http.Client) ClientOption {
	return &httpClientOption{client: client}
}
