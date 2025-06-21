package camunda

import (
	"context"
	"fmt"
	"path"
)

type Resource struct {
	client *Client
}

// TODO: change Resources from string to binary
type ResourceDeployRequest struct {
	Resources string `json:"resources"`
	TenantId  string `json:"tenantId"`
}

// TODO: change Deployments from string to binary
type ResourceDeployResponse struct {
	TenantId      string `json:"resources"`
	DeploymentKey string `json:"tenantId"`
	Deployments   string `json:"deployments"`
}

func (r Resource) Deploy(ctx context.Context, request ResourceDeployRequest) (*ResourceDeployResponse, error) {
	u := r.client.baseURL
	u.Path = path.Join(r.client.baseURL.Path, "deployments")
	fmt.Println(u.Path)
	return nil, nil
}
