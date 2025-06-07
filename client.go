package camunda

import (
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	httpClient *http.Client
	baseURL    url.URL

	Cluster Cluster
}

// Option represents a configuration option for the Client
type Option func(*Client) error

// WithTransport sets a custom transport
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) error {
		if cookieAuthTransport, ok := c.httpClient.Transport.(cookieAuth); ok {
			cookieAuthTransport.wrapped = transport
			c.httpClient.Transport = cookieAuthTransport
		} else {
			c.httpClient.Transport = transport
		}
		return nil
	}
}

// WithBaseURL sets the base URL for the client
func WithBaseURL(baseURL url.URL) Option {
	// Ensure path ends with /v2
	if !strings.HasSuffix(baseURL.Path, "v2") {
		baseURL.Path, _ = url.JoinPath(baseURL.Path, "v2")
	}

	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

// getTransport safely gets the transport from a client
func getTransport(client *http.Client) http.RoundTripper {
	if client.Transport == nil {
		return http.DefaultTransport
	}
	return client.Transport
}

// NewClient creates a new client with the given options
func NewClient(opts ...Option) (*Client, error) {
	client := &Client{
		httpClient: &http.Client{},
		baseURL: url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
			Path:   "v2",
		},
	}

	// Apply all opts
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return &Client{
		Cluster: Cluster{
			client: client,
		},
	}, nil
}
