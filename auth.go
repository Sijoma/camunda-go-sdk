package camunda

import (
	"context"
	"net/http"

	"golang.org/x/oauth2/clientcredentials"
)

type cookieAuth struct {
	wrapped http.RoundTripper
	value   string
}

func (m cookieAuth) RoundTrip(request *http.Request) (*http.Response, error) {
	cloned := request.Clone(request.Context())

	cloned.AddCookie(&http.Cookie{Name: "OPERATE-SESSION", Value: m.value})

	rt := m.wrapped
	if rt == nil {
		rt = http.DefaultTransport
	}
	return rt.RoundTrip(cloned)
}

// WithCookieAuth sets the cookie authentication value
func WithCookieAuth(value string) Option {
	return func(c *Client) {
		c.httpClient.Transport = cookieAuth{
			wrapped: getTransport(c.httpClient),
			value:   value,
		}
	}
}

func WithOAuth(clientID, clientSecret, tokenURL, audience string, scopes ...string) func(c *Client) {
	// OAuth2 config for client credentials flow
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       scopes, // Optional
		EndpointParams: map[string][]string{
			"audience": {audience},
		},
	}

	// Create a new HTTP client that automatically authenticates
	client := config.Client(context.Background())
	return func(c *Client) {
		c.httpClient = client
	}
}
