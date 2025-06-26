package internal

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/sijoma/camunda-go-sdk/orchestration"
)

type TestSuite struct {
	host               string
	gatewayEndpoint    string
	managementEndpoint string
	camundaContainer   testcontainers.Container
	t                  testing.TB
}

func NewTestSuite(t testing.TB, ctx context.Context) (*TestSuite, error) {
	ts := &TestSuite{t: t}
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ts.teardown(cleanupCtx)
	})

	err := ts.setupCamunda(ctx)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (ts *TestSuite) CamundaEndpoint() (string, error) {
	if ts.gatewayEndpoint == "" {
		return "", fmt.Errorf("gateway endpoint not initialized")
	}
	return ts.gatewayEndpoint, nil
}

func (ts *TestSuite) ManagementEndpoint() (string, error) {
	if ts.managementEndpoint == "" {
		return "", fmt.Errorf("management endpoint not initialized")
	}
	return ts.managementEndpoint, nil
}

func (ts *TestSuite) setupCamunda(ctx context.Context) error {
	req := testcontainers.ContainerRequest{
		Image:        "camunda/zeebe:latest",
		ExposedPorts: []string{"26500/tcp", "8080/tcp", "9600/tcp"},
		WaitingFor: wait.ForAll(
			&topologyWaitStrategy{},
			wait.ForLog("Tomcat started on port 8080 (http) with context path"),
			wait.ForLog("Tomcat started on port 9600"),
		).WithDeadline(2 * time.Minute),
	}

	camundaContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}
	ts.camundaContainer = camundaContainer

	ts.gatewayEndpoint, err = camundaContainer.PortEndpoint(ctx, "8080/tcp", "http")
	if err != nil {
		return err
	}

	ts.managementEndpoint, err = camundaContainer.PortEndpoint(ctx, "9600/tcp", "http")
	if err != nil {
		return err
	}

	return nil
}

func (ts *TestSuite) teardown(ctx context.Context) {
	if ts.camundaContainer != nil {
		if err := ts.camundaContainer.Terminate(ctx); err != nil {
			ts.t.Errorf("failed to terminate container: %v", err)
		}
	}
}

type topologyWaitStrategy struct{}

func (t topologyWaitStrategy) WaitUntilReady(ctx context.Context, target wait.StrategyTarget) error {
	maxAttempts := 30
	backoff := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for topology: %w", ctx.Err())
		default:
			host, err := target.Host(ctx)
			if err != nil {
				return err
			}

			zeebePort, err := nat.NewPort("", "8080")
			if err != nil {
				return err
			}

			mappedPort, err := target.MappedPort(ctx, zeebePort)
			if err != nil {
				return err
			}

			baseURL, err := url.Parse(fmt.Sprintf("http://%s:%d", host, mappedPort.Int()))
			if err != nil {
				return err
			}

			client, err := orchestration.NewClient(orchestration.WithBaseURL(*baseURL))
			if err != nil {
				return err
			}

			topology, err := client.Cluster.Topology(ctx)
			if err != nil {
				fmt.Printf("TestSetup HealthCheck: Attempt %d/%d: Failed to get topology: %v\n", attempt, maxAttempts, err)
			} else if topology.ClusterSize > 0 && len(topology.Brokers) > 0 {
				return nil // Success!
			}

			if attempt < maxAttempts {
				time.Sleep(backoff)
			}
		}
	}

	return fmt.Errorf("timeout waiting for cluster topology after %d attempts", maxAttempts)

}
