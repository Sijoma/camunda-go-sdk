package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	camunda "github.com/camunda/go-sdk"
)

//export ZEEBE_ADDRESS='54f3b165-c547-4fbd-98bd-60f0f72001da.ont-1.zeebe.camunda.io:443'
//export ZEEBE_CLIENT_ID='mdu54ZmcJ-CBcc.8JLiPEQqT-g8qXW40'
//export ZEEBE_CLIENT_SECRET='SXzgMonP5RI~l2U1naDjMKjm7nir1h1nD0cS8AVA0p-F-aY.EpHrb~nL68.asaLi'
//export ZEEBE_AUTHORIZATION_SERVER_URL='https://login.cloud.camunda.io/oauth/token'
//export ZEEBE_REST_ADDRESS='https://ont-1.zeebe.camunda.io/54f3b165-c547-4fbd-98bd-60f0f72001da'
//export ZEEBE_GRPC_ADDRESS='grpcs://54f3b165-c547-4fbd-98bd-60f0f72001da.ont-1.zeebe.camunda.io:443'
//export ZEEBE_TOKEN_AUDIENCE='zeebe.camunda.io'
//export CAMUNDA_CLUSTER_ID='54f3b165-c547-4fbd-98bd-60f0f72001da'
//export CAMUNDA_CLIENT_ID='mdu54ZmcJ-CBcc.8JLiPEQqT-g8qXW40'
//export CAMUNDA_CLIENT_SECRET='SXzgMonP5RI~l2U1naDjMKjm7nir1h1nD0cS8AVA0p-F-aY.EpHrb~nL68.asaLi'
//export CAMUNDA_CLUSTER_REGION='ont-1'
//export CAMUNDA_CREDENTIALS_SCOPES='Zeebe'
//export CAMUNDA_OAUTH_URL='https://login.cloud.camunda.io/oauth/token'
//export CAMUNDA_CLIENT_MODE='saas'
//export CAMUNDA_CLIENT_AUTH_CLIENTID='mdu54ZmcJ-CBcc.8JLiPEQqT-g8qXW40'
//export CAMUNDA_CLIENT_AUTH_CLIENTSECRET='SXzgMonP5RI~l2U1naDjMKjm7nir1h1nD0cS8AVA0p-F-aY.EpHrb~nL68.asaLi'
//export CAMUNDA_CLIENT_CLOUD_CLUSTERID='54f3b165-c547-4fbd-98bd-60f0f72001da'
//export CAMUNDA_CLIENT_CLOUD_REGION='ont-1'

func main() {
	ctx := context.Background()
	clientID := os.Getenv("CAMUNDA_CLIENT_ID")
	clientSecret := os.Getenv("CAMUNDA_CLIENT_SECRET")
	tokenURL := os.Getenv("CAMUNDA_OAUTH_URL")
	audience := os.Getenv("ZEEBE_TOKEN_AUDIENCE")
	scopes := os.Getenv("CAMUNDA_CREDENTIALS_SCOPES")
	baseURL, _ := url.Parse("http://localhost:8080")

	fmt.Println(clientID, clientSecret, tokenURL, audience, scopes, baseURL)

	c8 := camunda.NewClient(
		camunda.WithBaseURL(*baseURL),
		//camunda.WithOAuth(clientID, clientSecret, tokenURL, audience, scopes),
	)
	topology, err := c8.Cluster.Topology(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(topology)

}
