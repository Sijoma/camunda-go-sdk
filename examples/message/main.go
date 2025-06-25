package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/sijoma/camunda-go-sdk/orchestration"
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

	c8, err := orchestration.NewClient(
		orchestration.WithBaseURL(*baseURL),
		orchestration.WithOAuth(clientID, clientSecret, tokenURL, audience, scopes),
	)
	if err != nil {
		fmt.Println("failed creating client", err)
		return
	}

	payload := orchestration.MessagePublishRequest{
		Name:           "a name",
		CorrelationKey: "correlation-key",
		TimeToLive:     0,
		MessageId:      "no-id",
	}

	topology, err := c8.Message.Publish(ctx, payload)
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
