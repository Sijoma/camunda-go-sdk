package camunda

import (
	"context"
	"encoding/base64"
	"fmt"
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
	return func(c *Client) error {
		c.httpClient.Transport = cookieAuth{
			wrapped: getTransport(c.httpClient),
			value:   value,
		}
		return nil
	}
}

type basicAuth struct {
	wrapped http.RoundTripper
	value   string
}

func (m basicAuth) RoundTrip(request *http.Request) (*http.Response, error) {
	cloned := request.Clone(request.Context())
	cloned.Header.Set("Authorization", m.value)

	rt := m.wrapped
	if rt == nil {
		rt = http.DefaultTransport
	}
	return rt.RoundTrip(cloned)
}

func WithBasicAuth(username, password string) Option {
	return func(c *Client) error {
		if username == "" || password == "" {
			return fmt.Errorf("WithBasicAuth: username or password cannot be empty")
		}

		credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
		auth := "Basic " + credentials
		c.httpClient.Transport = basicAuth{
			wrapped: getTransport(c.httpClient),
			value:   auth,
		}
		return nil
	}
}

func WithOAuth(clientID, clientSecret, tokenURL, audience string, scopes ...string) func(c *Client) error {
	return func(c *Client) error {
		if clientID == "" || clientSecret == "" || tokenURL == "" || audience == "" {
			return fmt.Errorf("WithOAuth: clientID, clientSecret, tokenURL, audience cannot be empty")
		}
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
		c.httpClient = client
		return nil
	}
}
