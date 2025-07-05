package main

import (
	"context"
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

	payload := orchestration.ResourceDeployRequest{
		Resources: "test",
		TenantId:  "test",
	}

	resource, err := c8.Resource.Deploy(ctx, payload, "/Users/hamzamasood/Desktop/diagram_1.bpmn")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resource)
}
