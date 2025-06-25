package camunda_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sijoma/camunda-go-sdk"
	"github.com/sijoma/camunda-go-sdk/internal"
)

func Test_Camunda_Message(t *testing.T) {
	ctx := context.Background()
	suite, err := internal.NewTestSuite(t, ctx)
	if err != nil {
		t.Fatalf("failed to setup camunda: %v", err)
	}
	// Setup Client
	camundaURL, err := suite.CamundaEndpoint()
	require.NoError(t, err)
	baseURL, _ := url.Parse(camundaURL)
	c8, err := camunda.NewOrchestrationClusterClient(
		camunda.WithBaseURL(*baseURL),
	)
	require.NoError(t, err)

	t.Run("Publish Message", func(t *testing.T) {
		resp, err := c8.Message.Publish(t.Context(), camunda.MessagePublishRequest{
			Name:           "a name",
			CorrelationKey: "a correlation key",
			TimeToLive:     0,
			MessageId:      "message-id",
		})
		require.NoError(t, err)

		assert.NotEmptyf(t, resp.MessageKey, "message key should not be empty")
		assert.Equal(t, "<default>", resp.TenantId)
	})

	t.Run("Correlate Message", func(t *testing.T) {
		resp, err := c8.Message.Correlate(t.Context(), camunda.MessageCorrelateRequest{
			Name:           "a name",
			CorrelationKey: "a correlation key",
			Variables: map[string]interface{}{
				"foo": "bar",
			},
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Command 'CORRELATE' rejected with code 'NOT_FOUND'")
		assert.Nil(t, resp)
	})
}
