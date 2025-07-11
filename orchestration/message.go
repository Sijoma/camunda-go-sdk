package orchestration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

type Message struct {
	client *Client
}

type MessagePublishRequest struct {
	Name           string         `json:"name"`
	CorrelationKey string         `json:"correlationKey"`
	TimeToLive     int            `json:"timeToLive"`
	MessageId      string         `json:"messageId"`
	Variables      map[string]any `json:"variables"`
	TenantId       string         `json:"tenantId"`
}

type MessagePublishResponse struct {
	TenantId   string `json:"tenantId"`
	MessageKey string `json:"messageKey"`
}

func (m Message) Publish(ctx context.Context, request MessagePublishRequest) (*MessagePublishResponse, error) {
	u := m.client.baseURL
	u.Path = path.Join(m.client.baseURL.Path, "messages/publication")

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("publish marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := m.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		var publishResponse MessagePublishResponse
		err = json.NewDecoder(res.Body).Decode(&publishResponse)
		if err != nil {
			return nil, fmt.Errorf("publish: %w", err)
		}
		return &publishResponse, nil
	default:
		return nil, handleErrorResponse(res.Body, "publish")
	}
}

type MessageCorrelateRequest struct {
	Name           string `json:"name"`
	CorrelationKey string `json:"correlationKey"`
	Variables      any    `json:"variables"`
	TenantId       string `json:"tenantId"`
}

type MessageCorrelateResponse struct {
	TenantId           string `json:"tenantId"`
	MessageKey         string `json:"messageKey"`
	ProcessInstanceKey string `json:"processInstanceKey"`
}

func (m Message) Correlate(ctx context.Context, request MessageCorrelateRequest) (*MessageCorrelateResponse, error) {
	u := m.client.baseURL
	u.Path = path.Join(m.client.baseURL.Path, "messages/correlation")

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("correlate marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := m.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		var correlateResponse MessageCorrelateResponse
		err = json.NewDecoder(res.Body).Decode(&correlateResponse)
		if err != nil {
			return nil, fmt.Errorf("correlate: %w", err)
		}
		return &correlateResponse, nil
	default:
		return nil, handleErrorResponse(res.Body, "correlate")
	}
}
