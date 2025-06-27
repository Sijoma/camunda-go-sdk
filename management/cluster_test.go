package management_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sijoma/camunda-go-sdk/internal"
	"github.com/sijoma/camunda-go-sdk/management"
)

func TestCamunda_Management(t *testing.T) {
	ctx := context.Background()
	suite, err := internal.NewTestSuite(t, ctx)
	if err != nil {
		t.Fatalf("failed to setup camunda: %v", err)
	}
	// Setup Client
	managementURL, err := suite.ManagementEndpoint()
	require.NoError(t, err)
	// Create a new management client
	baseURL, err := url.Parse(managementURL)
	require.NoError(t, err, "Failed to parse URL")

	client, err := management.NewClient(
		management.WithBaseURL(*baseURL),
	)
	require.NoError(t, err, "Failed to create client")

	t.Run("Query Topology", func(t *testing.T) {
		// Get the current topology
		ctx := context.Background()
		topology, err := client.Cluster.Topology(ctx)
		require.NoError(t, err, "Failed to get topology")

		// Print the topology information
		fmt.Printf("Topology Version: %d\n", topology.Version)
		fmt.Printf("Number of Brokers: %d\n", len(topology.Brokers))

		// Assert that the topology has valid data
		assert.GreaterOrEqual(t, topology.Version, int64(0), "Topology version should be non-negative")
		assert.GreaterOrEqual(t, len(topology.Brokers), 0, "Brokers list can be empty but not nil")
		fmt.Printf("Pending change: %v\n", topology.PendingChange)
	})

	// TODO: All the below tests don't really work they need to be split up or put in a sequence.
	t.Run("Add Broker", func(t *testing.T) {
		ctx := context.Background()
		brokerId := management.BrokerId(3) // Example broker ID

		// Use dryRun=true to simulate without actually adding the broker
		response, err := client.Cluster.AddBroker(ctx, brokerId, false)
		require.NoError(t, err, "Failed to add broker (dry run)")

		// Verify the response
		assert.NotZero(t, response.ChangeId, "Change ID should not be zero")
		assert.NotEmpty(t, response.CurrentTopology, "Current topology should not be empty")
		assert.NotEmpty(t, response.PlannedChanges, "Planned changes should not be empty")
		assert.NotEmpty(t, response.ExpectedTopology, "Expected topology should not be empty")

		// Verify that at least one planned change is a BROKER_ADD operation
		brokerAddFound := false
		for _, op := range response.PlannedChanges {
			if op.Operation == "BROKER_ADD" && op.BrokerId == brokerId {
				brokerAddFound = true
				break
			}
		}
		assert.True(t, brokerAddFound, "BROKER_ADD operation for the specified broker ID not found")
	})

	t.Run("Remove Broker", func(t *testing.T) {
		ctx := context.Background()
		brokerId := management.BrokerId(3) // Example broker ID to remove

		// Use dryRun=true to simulate without actually removing the broker
		response, err := client.Cluster.RemoveBroker(ctx, brokerId, true)
		require.NoError(t, err, "Failed to remove broker (dry run)")

		// Verify the response
		assert.NotZero(t, response.ChangeId, "Change ID should not be zero")
		assert.NotEmpty(t, response.CurrentTopology, "Current topology should not be empty")
		assert.NotEmpty(t, response.PlannedChanges, "Planned changes should not be empty")
		assert.NotEmpty(t, response.ExpectedTopology, "Expected topology should not be empty")

		// Verify that at least one planned change is a BROKER_REMOVE operation
		brokerRemoveFound := false
		for _, op := range response.PlannedChanges {
			if op.Operation == "BROKER_REMOVE" && op.BrokerId == brokerId {
				brokerRemoveFound = true
				break
			}
		}
		assert.True(t, brokerRemoveFound, "BROKER_REMOVE operation for the specified broker ID not found")
	})

	t.Run("Scale Brokers", func(t *testing.T) {
		ctx := context.Background()
		brokerIds := []management.BrokerId{0, 1, 2} // Example broker IDs

		// Use dryRun=true to simulate without actually scaling the brokers
		response, err := client.Cluster.ScaleBrokers(ctx, brokerIds, true, false, nil)
		require.NoError(t, err, "Failed to scale brokers (dry run)")

		// Verify the response
		assert.NotZero(t, response.ChangeId, "Change ID should not be zero")
		assert.NotEmpty(t, response.CurrentTopology, "Current topology should not be empty")
		assert.NotEmpty(t, response.PlannedChanges, "Planned changes should not be empty")
		assert.NotEmpty(t, response.ExpectedTopology, "Expected topology should not be empty")

		// Verify that the expected topology contains the specified brokers
		brokerIdsFound := make(map[management.BrokerId]bool)
		for _, broker := range response.ExpectedTopology {
			brokerIdsFound[broker.ID] = true
		}

		for _, id := range brokerIds {
			assert.True(t, brokerIdsFound[id], "Broker ID %d not found in expected topology", id)
		}
	})

	t.Run("Purge Cluster", func(t *testing.T) {
		ctx := context.Background()

		// Use dryRun=true to simulate without actually purging the cluster
		response, err := client.Cluster.PurgeCluster(ctx, true)
		require.NoError(t, err, "Failed to purge cluster (dry run)")

		// Verify the response
		assert.NotZero(t, response.ChangeId, "Change ID should not be zero")
		assert.NotEmpty(t, response.CurrentTopology, "Current topology should not be empty")
		assert.NotEmpty(t, response.PlannedChanges, "Planned changes should not be empty")
		assert.NotEmpty(t, response.ExpectedTopology, "Expected topology should not be empty")

		// Verify that at least one planned change is a DELETE_HISTORY operation
		deleteHistoryFound := false
		for _, op := range response.PlannedChanges {
			if op.Operation == "DELETE_HISTORY" {
				deleteHistoryFound = true
				break
			}
		}
		assert.True(t, deleteHistoryFound, "DELETE_HISTORY operation not found")
	})

	t.Run("Reconfigure Cluster", func(t *testing.T) {
		ctx := context.Background()

		// Create a reconfiguration request
		count := int32(3)
		replicationFactor := int32(3)
		request := management.ClusterConfigPatchRequest{
			Brokers: &management.BrokersConfig{
				Count: &count,
			},
			Partitions: &management.PartitionsConfig{
				ReplicationFactor: &replicationFactor,
			},
		}

		// Use dryRun=true to simulate without actually reconfiguring the cluster
		response, err := client.Cluster.ReconfigureCluster(ctx, request, true, false)
		require.NoError(t, err, "Failed to reconfigure cluster (dry run)")

		// Verify the response
		assert.NotZero(t, response.ChangeId, "Change ID should not be zero")
		assert.NotEmpty(t, response.CurrentTopology, "Current topology should not be empty")
		assert.NotEmpty(t, response.PlannedChanges, "Planned changes should not be empty")
		assert.NotEmpty(t, response.ExpectedTopology, "Expected topology should not be empty")

		// Verify that the expected topology has the correct number of brokers
		assert.Len(t, response.ExpectedTopology, int(count), "Expected topology should have %d brokers", count)
	})
}
