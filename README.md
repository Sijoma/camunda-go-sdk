# Camunda Go SDK

A Go client SDK for Camunda 8. 🚧⚠️🏗️🔨

So far only supports the REST API and one endpoint. 

## Installation

```bash
go get github.com/sijoma/camunda-go-sdk
```

Topology Request:
```go
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
```

Output
```json
{
  "brokers": [
    {
      "nodeId": 0,
      "host": "zeebe-0.zeebe-broker-service.cc17ab23-5cc1-47ef-9e01-9957c6e10911-zeebe.svc.cluster.local",
      "port": 26501,
      "partitions": [
        {
          "partitionId": 1,
          "role": "leader",
          "health": "healthy"
        },
        {
          "partitionId": 3,
          "role": "follower",
          "health": "healthy"
        },
        {
          "partitionId": 2,
          "role": "follower",
          "health": "healthy"
        }
      ],
      "version": "8.7.2"
    },
    {
      "nodeId": 1,
      "host": "zeebe-1.zeebe-broker-service.cc17ab23-5cc1-47ef-9e01-9957c6e10911-zeebe.svc.cluster.local",
      "port": 26501,
      "partitions": [
        {
          "partitionId": 1,
          "role": "follower",
          "health": "healthy"
        },
        {
          "partitionId": 3,
          "role": "follower",
          "health": "healthy"
        },
        {
          "partitionId": 2,
          "role": "leader",
          "health": "healthy"
        }
      ],
      "version": "8.7.2"
    },
    {
      "nodeId": 2,
      "host": "zeebe-2.zeebe-broker-service.cc17ab23-5cc1-47ef-9e01-9957c6e10911-zeebe.svc.cluster.local",
      "port": 26501,
      "partitions": [
        {
          "partitionId": 1,
          "role": "follower",
          "health": "healthy"
        },
        {
          "partitionId": 3,
          "role": "leader",
          "health": "healthy"
        },
        {
          "partitionId": 2,
          "role": "follower",
          "health": "healthy"
        }
      ],
      "version": "8.7.2"
    }
  ],
  "clusterSize": 3,
  "partitionsCount": 3,
  "replicationFactor": 3,
  "gatewayVersion": "8.7.2",
  "lastCompletedChangeId": ""
}

```
