package camunda

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"path"
)

type Cluster struct {
	client *Client
}

func (c Cluster) Topology(ctx context.Context) (*TopologyResponse, error) {
	u := c.client.baseURL
	u.Path = path.Join(c.client.baseURL.Path, "topology")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 299 {
		dump, _ := httputil.DumpResponse(res, true)
		return nil, fmt.Errorf("received status code %d: %s", res.StatusCode, dump)
	}

	var topology TopologyResponse
	err = json.NewDecoder(res.Body).Decode(
		&topology,
	)
	if err != nil {
		return nil, err
	}

	return &topology, nil
}

type TopologyResponse struct {
	Brokers []struct {
		NodeId     int    `json:"nodeId"`
		Host       string `json:"host"`
		Port       int    `json:"port"`
		Partitions []struct {
			PartitionId int    `json:"partitionId"`
			Role        string `json:"role"`
			Health      string `json:"health"`
		} `json:"partitions"`
		Version string `json:"version"`
	} `json:"brokers"`
	ClusterSize           int    `json:"clusterSize"`
	PartitionsCount       int    `json:"partitionsCount"`
	ReplicationFactor     int    `json:"replicationFactor"`
	GatewayVersion        string `json:"gatewayVersion"`
	LastCompletedChangeId string `json:"lastCompletedChangeId"`
}
