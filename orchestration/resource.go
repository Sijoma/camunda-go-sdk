package orchestration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type Resource struct {
	client *Client
}

type ResourceGetRequest struct {
	ResourceKey string `json:"resourceKey"`
}

type ResourceGetResponse struct {
	ResourceName string `json:"resourceName"`
	Version      int    `json:"version"`
	VersionTag   string `json:"versionTag"`
	ResourceId   string `json:"resourceId"`
	TenantId     string `json:"tenantId"`
	ResourceKey  string `json:"resourceKey"`
}

func (r Resource) Get(ctx context.Context, request ResourceGetRequest) (*ResourceGetResponse, error) {
	u := r.client.baseURL
	u.Path = path.Join(r.client.baseURL.Path, "resources", request.ResourceKey)
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := r.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		var getResponse ResourceGetResponse
		err = json.NewDecoder(res.Body).Decode(&getResponse)
		if err != nil {
			return nil, fmt.Errorf("publish: %w", err)
		}
		return &getResponse, nil
	default:
		var errorResponse ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errorResponse)
		if err != nil {
			return nil, fmt.Errorf("couldn't decode response: %w", err)
		}
		return nil, fmt.Errorf("get resource error title: %s\nget resource error detail %s", errorResponse.Title, errorResponse.Detail)
	}
}
