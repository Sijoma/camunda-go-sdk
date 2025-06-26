package management

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"path"
	"strconv"
)

// Cluster provides methods for interacting with the Cluster Topology Management API
type Cluster struct {
	client *Client
}

// BrokerId represents the ID of a broker
type BrokerId int32

// PartitionId represents the ID of a partition
type PartitionId int32

// ChangeId represents the ID of a topology change operation
type ChangeId int64

// TopologyVersion represents the version of the topology
type TopologyVersion int64

// BrokerState represents the state of a broker
type BrokerState struct {
	ID            BrokerId         `json:"id"`
	State         string           `json:"state"`
	Version       int64            `json:"version"`
	LastUpdatedAt string           `json:"lastUpdatedAt"`
	Partitions    []PartitionState `json:"partitions"`
}

// PartitionState represents the state of a partition
type PartitionState struct {
	ID       PartitionId     `json:"id"`
	State    string          `json:"state"`
	Priority int32           `json:"priority"`
	Config   PartitionConfig `json:"config,omitempty"`
}

// PartitionConfig represents the configuration of a partition
type PartitionConfig struct {
	Exporting ExportingConfig `json:"exporting,omitempty"`
}

// ExportingConfig represents the exporting configuration
type ExportingConfig struct {
	Exporters []ExporterConfig `json:"exporters,omitempty"`
}

// ExporterConfig represents the configuration of an exporter
type ExporterConfig struct {
	ID    string `json:"id"`
	State string `json:"state"`
}

// Operation represents a topology change operation
type Operation struct {
	Operation   string      `json:"operation"`
	BrokerId    BrokerId    `json:"brokerId,omitempty"`
	PartitionId PartitionId `json:"partitionId,omitempty"`
	Priority    int32       `json:"priority,omitempty"`
	Brokers     []BrokerId  `json:"brokers,omitempty"`
	ExporterId  string      `json:"exporterId,omitempty"`
}

// PlannedOperationsResponse represents the response for operations that plan changes
type PlannedOperationsResponse struct {
	ChangeId         ChangeId      `json:"changeId"`
	CurrentTopology  []BrokerState `json:"currentTopology"`
	PlannedChanges   []Operation   `json:"plannedChanges"`
	ExpectedTopology []BrokerState `json:"expectedTopology"`
}

// TopologyResponse represents the response for getting the current topology
type TopologyResponse struct {
	Version       TopologyVersion `json:"version"`
	Brokers       []BrokerState   `json:"brokers"`
	LastChange    CompletedChange `json:"lastChange,omitempty"`
	PendingChange TopologyChange  `json:"pendingChange,omitempty"`
	Routing       RoutingState    `json:"routing,omitempty"`
}

// CompletedChange represents a completed topology change
type CompletedChange struct {
	ID          ChangeId `json:"id"`
	Status      string   `json:"status"`
	StartedAt   string   `json:"startedAt"`
	CompletedAt string   `json:"completedAt"`
}

// TopologyChange represents a topology change in progress
type TopologyChange struct {
	ID              ChangeId    `json:"id"`
	Status          string      `json:"status"`
	StartedAt       string      `json:"startedAt"`
	CompletedAt     string      `json:"completedAt,omitempty"`
	InternalVersion int64       `json:"internalVersion"`
	Completed       []Operation `json:"completed,omitempty"`
	Pending         []Operation `json:"pending,omitempty"`
}

// RoutingState represents the routing state
type RoutingState struct {
	Version            int64              `json:"version"`
	RequestHandling    RequestHandling    `json:"requestHandling"`
	MessageCorrelation MessageCorrelation `json:"messageCorrelation"`
}

// RequestHandling represents the request handling strategy
type RequestHandling struct {
	Strategy string `json:"strategy"`
	// Additional fields based on strategy
}

// MessageCorrelation represents the message correlation strategy
type MessageCorrelation struct {
	Strategy string `json:"strategy"`
	// Additional fields based on strategy
}

// ScaleRequest represents the request body for scaling brokers
type ScaleRequest []BrokerId

// ClusterConfigPatchRequest represents the request body for patching cluster configuration
type ClusterConfigPatchRequest struct {
	Brokers    *BrokersConfig    `json:"brokers,omitempty"`
	Partitions *PartitionsConfig `json:"partitions,omitempty"`
}

// BrokersConfig represents the brokers configuration for a patch request
type BrokersConfig struct {
	Add    []BrokerId `json:"add,omitempty"`
	Remove []BrokerId `json:"remove,omitempty"`
	Count  *int32     `json:"count,omitempty"`
}

// PartitionsConfig represents the partitions configuration for a patch request
type PartitionsConfig struct {
	Count             *int32 `json:"count,omitempty"`
	ReplicationFactor *int32 `json:"replicationFactor,omitempty"`
}

// Topology gets the current topology of the cluster
func (c Cluster) Topology(ctx context.Context) (*TopologyResponse, error) {
	u := c.client.baseURL

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

	if res.StatusCode >= 300 {
		dump, _ := httputil.DumpResponse(res, true)
		return nil, fmt.Errorf("received status code %d: %s", res.StatusCode, dump)
	}

	var topology TopologyResponse
	err = json.NewDecoder(res.Body).Decode(&topology)
	if err != nil {
		return nil, err
	}

	return &topology, nil
}

// AddBroker adds a broker with the given brokerId to the cluster
func (c Cluster) AddBroker(ctx context.Context, brokerId BrokerId, dryRun bool) (*PlannedOperationsResponse, error) {
	u := c.client.baseURL
	u.Path = path.Join(c.client.baseURL.Path, "brokers", strconv.Itoa(int(brokerId)))

	// Add query parameters
	q := u.Query()
	if dryRun {
		q.Set("dryRun", "true")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		dump, _ := httputil.DumpResponse(res, true)
		return nil, fmt.Errorf("received status code %d: %s", res.StatusCode, dump)
	}

	var response PlannedOperationsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RemoveBroker removes a broker with the given brokerId from the cluster
func (c Cluster) RemoveBroker(ctx context.Context, brokerId BrokerId, dryRun bool) (*PlannedOperationsResponse, error) {
	u := c.client.baseURL
	u.Path = path.Join(c.client.baseURL.Path, "brokers", strconv.Itoa(int(brokerId)))

	// Add query parameters
	q := u.Query()
	if dryRun {
		q.Set("dryRun", "true")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		dump, _ := httputil.DumpResponse(res, true)
		return nil, fmt.Errorf("received status code %d: %s", res.StatusCode, dump)
	}

	var response PlannedOperationsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ScaleBrokers reconfigures the cluster with the given brokers
func (c Cluster) ScaleBrokers(ctx context.Context, brokerIds []BrokerId, dryRun bool, force bool, replicationFactor *int32) (*PlannedOperationsResponse, error) {
	u := c.client.baseURL
	u.Path = path.Join(c.client.baseURL.Path, "brokers")

	// Add query parameters
	q := u.Query()
	if dryRun {
		q.Set("dryRun", "true")
	}
	if force {
		q.Set("force", "true")
	}
	if replicationFactor != nil {
		q.Set("replicationFactor", strconv.Itoa(int(*replicationFactor)))
	}
	u.RawQuery = q.Encode()

	// Create request body
	body, err := json.Marshal(brokerIds)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		dump, _ := httputil.DumpResponse(res, true)
		return nil, fmt.Errorf("received status code %d: %s", res.StatusCode, dump)
	}

	var response PlannedOperationsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PurgeCluster purges data from the cluster
func (c Cluster) PurgeCluster(ctx context.Context, dryRun bool) (*PlannedOperationsResponse, error) {
	u := c.client.baseURL
	u.Path = path.Join(c.client.baseURL.Path, "purge")

	// Add query parameters
	q := u.Query()
	if dryRun {
		q.Set("dryRun", "true")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		dump, _ := httputil.DumpResponse(res, true)
		return nil, fmt.Errorf("received status code %d: %s", res.StatusCode, dump)
	}

	var response PlannedOperationsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ReconfigureCluster reconfigures the cluster by adding or removing brokers, adding more partitions, or changing the replication factor
func (c Cluster) ReconfigureCluster(ctx context.Context, request ClusterConfigPatchRequest, dryRun bool, force bool) (*PlannedOperationsResponse, error) {
	u := c.client.baseURL

	// Add query parameters
	q := u.Query()
	if dryRun {
		q.Set("dryRun", "true")
	}
	if force {
		q.Set("force", "true")
	}
	u.RawQuery = q.Encode()

	// Create request body
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		dump, _ := httputil.DumpResponse(res, true)
		return nil, fmt.Errorf("received status code %d: %s", res.StatusCode, dump)
	}

	var response PlannedOperationsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
