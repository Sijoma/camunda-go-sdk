package camunda

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

type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
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
		var errorResponse ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errorResponse)
		if err != nil {
			return nil, fmt.Errorf("publish: %w", err)
		}
		return nil, fmt.Errorf("publish: %s", errorResponse.Detail)
	}
}
