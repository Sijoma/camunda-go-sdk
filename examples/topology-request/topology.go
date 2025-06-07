package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/sijoma/camunda-go-sdk"
)

func main() {
	ctx := context.Background()
	clientID := os.Getenv("CAMUNDA_CLIENT_ID")
	clientSecret := os.Getenv("CAMUNDA_CLIENT_SECRET")
	tokenURL := os.Getenv("CAMUNDA_OAUTH_URL")
	audience := os.Getenv("ZEEBE_TOKEN_AUDIENCE")
	scopes := os.Getenv("CAMUNDA_CREDENTIALS_SCOPES")
	zeebeAddress := os.Getenv("ZEEBE_REST_ADDRESS")
	baseURL, err := url.Parse(zeebeAddress)
	if err != nil {
		panic(err)
	}

	fmt.Println("EXAMPLE: Config", clientID, clientSecret, tokenURL, audience, scopes, baseURL)

	c8 := camunda.NewClient(
		camunda.WithBaseURL(*baseURL),
		camunda.WithOAuth(clientID, clientSecret, tokenURL, audience, scopes),
	)
	topology, err := c8.Cluster.Topology(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	pretty, err := json.MarshalIndent(topology, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(pretty))
}
