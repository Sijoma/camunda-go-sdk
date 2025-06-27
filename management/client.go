package management

import (
	"net/http"
	"net/url"
	"strings"
)

// Client is the main client for interacting with the Camunda Management API
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
		c.httpClient.Transport = transport
		return nil
	}
}

// WithBaseURL sets the base URL for the management API
// it makes sure that the actuator/cluster path is present
func WithBaseURL(baseURL url.URL) Option {
	// Ensure path ends with /actuator/cluster
	if !strings.HasSuffix(baseURL.Path, "actuator/cluster") {
		baseURL.Path, _ = url.JoinPath(baseURL.Path, "actuator/cluster")
	}

	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

// NewClient creates a new client with the given options
func NewClient(opts ...Option) (*Client, error) {
	client := &Client{
		httpClient: &http.Client{},
		baseURL: url.URL{
			Scheme: "http",
			Host:   "localhost:9600",
			Path:   "actuator/cluster",
		},
	}

	// Apply all opts
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return &Client{
		httpClient: client.httpClient,
		baseURL:    client.baseURL,
		Cluster:    Cluster{client: client},
	}, nil
}
