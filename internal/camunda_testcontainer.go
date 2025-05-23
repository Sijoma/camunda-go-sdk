package internal

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestSuite struct {
	host             string
	gatewayPort      nat.Port
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
	if ts.host == "" || ts.gatewayPort.Int() == 0 {
		return "", errors.New("host or gatewayPort not initialized")
	}
	return "http://" + ts.host + ":" + ts.gatewayPort.Port(), nil
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

	host, err := camundaContainer.Host(ctx)
	if err != nil {
		return err
	}
	port, err := camundaContainer.MappedPort(ctx, "8080")
	if err != nil {
		return err
	}

	ts.host = host
	ts.gatewayPort = port
	return nil
}

func (ts *TestSuite) teardown(ctx context.Context) {
	if ts.camundaContainer != nil {
		if err := ts.camundaContainer.Terminate(ctx); err != nil {
			ts.t.Errorf("failed to terminate container: %v", err)
		}
	}
}
