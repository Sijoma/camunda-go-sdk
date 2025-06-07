package internal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestSuite struct {
	host             string
	gatewayEndpoint  string
	camundaContainer testcontainers.Container
	t                testing.TB
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

func (ts *TestSuite) setupCamunda(ctx context.Context) error {
	req := testcontainers.ContainerRequest{
		Image:        "camunda/zeebe:latest",
		ExposedPorts: []string{"26500/tcp", "8080/tcp"},
		WaitingFor: wait.ForAll(
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
	return nil
}

func (ts *TestSuite) teardown(ctx context.Context) {
	if ts.camundaContainer != nil {
		if err := ts.camundaContainer.Terminate(ctx); err != nil {
			ts.t.Errorf("failed to terminate container: %v", err)
		}
	}
}
