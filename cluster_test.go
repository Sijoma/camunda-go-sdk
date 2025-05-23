package camunda_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/camunda/go-sdk"
	"github.com/camunda/go-sdk/internal"
)

func TestWithCamunda(t *testing.T) {
	ctx := context.Background()
	suite, err := internal.NewTestSuite(t, ctx)
	if err != nil {
		t.Fatalf("failed to setup camunda: %v", err)
	}

	t.Run("Query Topology", func(t *testing.T) {
		camundaURL, err := suite.CamundaEndpoint()
		require.NoError(t, err)
		baseURL, _ := url.Parse(camundaURL)
		c8 := camunda.NewClient(
			camunda.WithBaseURL(*baseURL),
		)
		topology, err := c8.Cluster.Topology(t.Context())
		require.NoError(t, err)

		assert.Equal(t, 1, topology.PartitionsCount)
	})
}
