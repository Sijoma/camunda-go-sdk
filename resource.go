package camunda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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

func (r Resource) Deploy(ctx context.Context, request ResourceDeployRequest, filePath string) (*ResourceDeployResponse, error) {
	u := r.client.baseURL
	u.Path = path.Join(r.client.baseURL.Path, "deployments")

	buf := &bytes.Buffer{}

	mpw := multipart.NewWriter(buf)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	fWriter, err := mpw.CreateFormFile("resources", filePath)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fWriter, f)
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()

	// Close the multipart writer before creating the request
	err = mpw.Close()
	if err != nil {
		return nil, err
	}

	// set up the request

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mpw.FormDataContentType()) // detect the form data content type

	resp, err := r.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(u.String())
	fmt.Println(body)
	fmt.Println(resp.StatusCode)

	//---
	fmt.Println(u.Path)
	return nil, nil
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
	fmt.Println(u.Path)
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("publish marshal: %w", err)
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
			return nil, fmt.Errorf("couldn't get error: %w", err)
		}
		return nil, fmt.Errorf("get resource error title: %s\nget resource error detail %s", errorResponse.Title, errorResponse.Detail)
	}
}
