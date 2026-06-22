package client

import "context"

// Client defines behavior expected by commands.
type Client interface {
	Ping(ctx context.Context) error
}

// APIClient is a placeholder production implementation.
type APIClient struct {
	endpoint string
	token    string
}

func New(endpoint, token string) *APIClient {
	return &APIClient{endpoint: endpoint, token: token}
}

func (c *APIClient) Ping(ctx context.Context) error {
	_ = ctx
	return nil
}
