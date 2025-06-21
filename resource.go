package camunda

import (
	"context"
	"fmt"
)

type Resource struct {
	client *Client
}

// TODO: change Resources from string to binary
type ResourceDeployRequest struct {
	Resources string `json:"resources"`
	TeantnId  string `json:"tenantId"`
}

// TODO: change Deployments from string to binary
type ResourceDeployResponse struct {
	TenantId      string `json:"resources"`
	DeploymentKey string `json:"tenantId"`
	Deployments   string `json:"deployments"`
}

func (r Resource) Deploy(ctx context.Context, request ResourceDeployRequest) (*ResourceDeployResponse, error) {
	u := r.client.baseURL
	fmt.Println(u)
	return nil, nil
}
