package orchestration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type Job struct {
	client *Client
}

type JobActivateRequest struct {
	Type              string   `json:"type"`
	Worker            string   `json:"worker"`
	Timeout           int64    `json:"timeout"`
	MaxJobsToActivate int32    `json:"maxJobsToActivate"`
	FetchVariable     []string `json:"fetchVariable,omitempty"`
	TenantIds         []string `json:"tenantIds,omitempty"`
	RequestTimeout    int64    `json:"requestTimeout,omitempty"`
}

type ActivatedJob struct {
	Key                      string            `json:"jobKey"`
	Type                     string            `json:"type"`
	ProcessInstanceKey       string            `json:"processInstanceKey"`
	ProcessDefinitionId      string            `json:"processDefinitionId"`
	ProcessDefinitionVersion int32             `json:"processDefinitionVersion"`
	ProcessDefinitionKey     string            `json:"processDefinitionKey"`
	ElementId                string            `json:"elementId"`
	ElementInstanceKey       string            `json:"elementInstanceKey"`
	CustomHeaders            map[string]string `json:"customHeaders"`
	Worker                   string            `json:"worker"`
	Retries                  int32             `json:"retries"`
	Deadline                 int64             `json:"deadline"`
	Variables                map[string]any    `json:"variables"`
	TenantId                 string            `json:"tenantId"`
}

type JobActivateResponse struct {
	Jobs []ActivatedJob `json:"jobs"`
}

func (j Job) Activate(ctx context.Context, request JobActivateRequest) ([]ActivatedJob, error) {
	u := j.client.baseURL
	u.Path = path.Join(j.client.baseURL.Path, "jobs/activation")

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("activate jobs marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := j.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		var response JobActivateResponse
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			return nil, fmt.Errorf("activate jobs decode: %w", err)
		}
		return response.Jobs, nil
	default:
		return nil, handleErrorResponse(res.Body, "activate jobs")
	}
}

type JobCompleteRequest struct {
	Variables map[string]any `json:"variables,omitempty"`
	Result    any            `json:"result,omitempty"`
}

func (j Job) Complete(ctx context.Context, jobKey string, request JobCompleteRequest) error {
	u := j.client.baseURL
	u.Path = path.Join(j.client.baseURL.Path, "jobs", jobKey, "completion")

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("complete job marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := j.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		return handleErrorResponse(res.Body, "complete job")
	}

	return nil
}

type JobFailRequest struct {
	Retries      int32          `json:"retries"`
	ErrorMessage string         `json:"errorMessage,omitempty"`
	RetryBackOff int64          `json:"retryBackOff,omitempty"`
	Variables    map[string]any `json:"variables,omitempty"`
}

func (j Job) Fail(ctx context.Context, jobKey string, request JobFailRequest) error {
	u := j.client.baseURL
	u.Path = path.Join(j.client.baseURL.Path, "jobs", jobKey, "failure")

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("fail job marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := j.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		return handleErrorResponse(res.Body, "fail job")
	}

	return nil
}

type JobErrorRequest struct {
	ErrorCode    string         `json:"errorCode"`
	ErrorMessage string         `json:"errorMessage,omitempty"`
	Variables    map[string]any `json:"variables,omitempty"`
}

func (j Job) Error(ctx context.Context, jobKey string, request JobErrorRequest) error {
	u := j.client.baseURL
	u.Path = path.Join(j.client.baseURL.Path, "jobs", jobKey, "error")

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("job error marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := j.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		return handleErrorResponse(res.Body, "job error")
	}

	return nil
}

type JobUpdateRequest struct {
	Changeset JobChangeset `json:"changeset"`
}

type JobChangeset struct {
	Retries *int32 `json:"retries,omitempty"`
	Timeout *int64 `json:"timeout,omitempty"`
}

func (j Job) UpdateRetries(ctx context.Context, jobKey string, request JobUpdateRequest) error {
	u := j.client.baseURL
	u.Path = path.Join(j.client.baseURL.Path, "jobs", jobKey)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("update job retries marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := j.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		return handleErrorResponse(res.Body, "update job retries")
	}

	return nil
}
