package camunda_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/camunda/go-sdk"
)

func TestWithCamunda(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "camunda/zeebe:latest",
		ExposedPorts: []string{"26500/tcp", "8080/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Tomcat started on port 8080 (http) with context path"),
			wait.ForLog("Tomcat started on port 9600"),
		),
	}
	camundaContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	require.NoError(t, err)

	host, err := camundaContainer.Host(t.Context())
	require.NoError(t, err)
	port, err := camundaContainer.MappedPort(t.Context(), "8080")
	require.NoError(t, err)

	// Test begins here
	baseURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", host, port.Int()))
	c8 := camunda.NewClient(
		camunda.WithBaseURL(*baseURL),
	)
	topology, err := c8.Cluster.Topology(t.Context())
	require.NoError(t, err)

	assert.Equal(t, 1, topology.PartitionsCount)

	// Test ends here

	testcontainers.CleanupContainer(t, camundaContainer)
	require.NoError(t, err)
}
