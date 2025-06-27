package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/sijoma/camunda-go-sdk/management"
)

func main() {
	// Create a new management client
	baseURL, err := url.Parse("http://localhost:9600")
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	client, err := management.NewClient(
		management.WithBaseURL(*baseURL),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get the current topology
	ctx := context.Background()
	topology, err := client.Cluster.Topology(ctx)
	if err != nil {
		log.Fatalf("Failed to get topology: %v", err)
	}

	// Print the topology information
	fmt.Printf("Topology Version: %d\n", topology.Version)
	fmt.Printf("Number of Brokers: %d\n", len(topology.Brokers))

	for _, broker := range topology.Brokers {
		fmt.Printf("\nBroker ID: %d\n", broker.ID)
		fmt.Printf("  State: %s\n", broker.State)
		fmt.Printf("  Version: %d\n", broker.Version)
		fmt.Printf("  Last Updated: %s\n", broker.LastUpdatedAt)
		fmt.Printf("  Partitions: %d\n", len(broker.Partitions))

		for _, partition := range broker.Partitions {
			fmt.Printf("    Partition ID: %d\n", partition.ID)
			fmt.Printf("      State: %s\n", partition.State)
			fmt.Printf("      Priority: %d\n", partition.Priority)
		}
	}

	// Example of adding a broker
	brokerId := management.BrokerId(3)                             // Example broker ID
	response, err := client.Cluster.AddBroker(ctx, brokerId, true) // Using dryRun=true to simulate
	if err != nil {
		log.Fatalf("Failed to add broker: %v", err)
	}
	fmt.Printf("\nChange ID: %d\n", response.ChangeId)
	fmt.Printf("Planned Changes: %d\n", len(response.PlannedChanges))
	for i, op := range response.PlannedChanges {
		fmt.Printf("  Operation %d: %s\n", i+1, op.Operation)
	}

	// Example of scaling brokers
	brokerIds := []management.BrokerId{0, 1, 2}                                   // Example broker IDs
	response, err = client.Cluster.ScaleBrokers(ctx, brokerIds, true, false, nil) // Using dryRun=true to simulate
	if err != nil {
		log.Fatalf("Failed to scale brokers: %v", err)
	}
	fmt.Printf("\nChange ID: %d\n", response.ChangeId)
	fmt.Printf("Planned Changes: %d\n", len(response.PlannedChanges))
	for i, op := range response.PlannedChanges {
		fmt.Printf("  Operation %d: %s\n", i+1, op.Operation)
	}

	// Example of reconfiguring the cluster
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
	response, err = client.Cluster.ReconfigureCluster(ctx, request, true, false) // Using dryRun=true to simulate
	if err != nil {
		log.Fatalf("Failed to reconfigure cluster: %v", err)
	}
	fmt.Printf("\nChange ID: %d\n", response.ChangeId)
	fmt.Printf("Planned Changes: %d\n", len(response.PlannedChanges))
	for i, op := range response.PlannedChanges {
		fmt.Printf("  Operation %d: %s\n", i+1, op.Operation)
	}
}
