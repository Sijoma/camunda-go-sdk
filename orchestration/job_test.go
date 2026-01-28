package orchestration_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sijoma/camunda-go-sdk/internal"
	"github.com/sijoma/camunda-go-sdk/orchestration"
)

func Test_Camunda_Job(t *testing.T) {
	ctx := context.Background()
	suite, err := internal.NewTestSuite(t, ctx)
	if err != nil {
		t.Fatalf("failed to setup camunda: %v", err)
	}

	camundaURL, err := suite.CamundaEndpoint()
	require.NoError(t, err)
	baseURL, _ := url.Parse(camundaURL)
	c8, err := orchestration.NewClient(
		orchestration.WithBaseURL(*baseURL),
	)
	require.NoError(t, err)

	t.Run("Activate Jobs - Empty", func(t *testing.T) {
		jobs, err := c8.Job.Activate(ctx, orchestration.JobActivateRequest{
			Type:              "non-existent",
			Worker:            "test-worker",
			Timeout:           1000,
			MaxJobsToActivate: 1,
		})
		require.NoError(t, err)
		assert.Empty(t, jobs)
	})

	t.Run("Complete Job - Invalid Key", func(t *testing.T) {
		err := c8.Job.Complete(ctx, "0", orchestration.JobCompleteRequest{})
		require.Error(t, err)
	})

	t.Run("Fail Job - Invalid Key", func(t *testing.T) {
		err := c8.Job.Fail(ctx, "0", orchestration.JobFailRequest{
			Retries: 0,
		})
		require.Error(t, err)
	})

	t.Run("Job Error - Invalid Key", func(t *testing.T) {
		err := c8.Job.Error(ctx, "0", orchestration.JobErrorRequest{
			ErrorCode: "test-error",
		})
		require.Error(t, err)
	})

	t.Run("Update Job Retries - Invalid Key", func(t *testing.T) {
		retries := int32(3)
		err := c8.Job.UpdateRetries(ctx, "0", orchestration.JobUpdateRequest{
			Changeset: orchestration.JobChangeset{
				Retries: &retries,
			},
		})
		require.Error(t, err)
	})

	t.Run("Worker - Poll and Stop", func(t *testing.T) {
		workerCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		worker := c8.NewJobWorker("test-type", func(client *orchestration.Client, job orchestration.ActivatedJob) {
			// should not be called
		}, orchestration.WithPollInterval(100*time.Millisecond))

		err := worker.Run(workerCtx)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
