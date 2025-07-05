package orchestration

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handleErrorResponse(t *testing.T) {
	t.Run("successful decoding", func(t *testing.T) {
		// Create a valid JSON response
		jsonBody := `{
			"type": "error-type",
			"title": "Error Title",
			"status": 400,
			"detail": "Error details",
			"instance": "error-instance"
		}`
		body := io.NopCloser(bytes.NewReader([]byte(jsonBody)))

		// Call the function with a prefix
		err := handleErrorResponse(body, "test-prefix")

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "test-prefix")
		assert.Contains(t, err.Error(), "error-type")
		assert.Contains(t, err.Error(), "400")
		assert.Contains(t, err.Error(), "Error details")
		assert.Contains(t, err.Error(), "error-instance")

		// Verify the error type
		var errorResponse ErrorResponse
		assert.True(t, errors.As(err, &errorResponse))
		assert.Equal(t, "error-type", errorResponse.Type)
		assert.Equal(t, "Error Title", errorResponse.Title)
		assert.Equal(t, 400, errorResponse.Status)
		assert.Equal(t, "Error details", errorResponse.Detail)
		assert.Equal(t, "error-instance", errorResponse.Instance)
	})

	t.Run("failed decoding", func(t *testing.T) {
		// Create an invalid JSON response
		invalidJSON := `{ "type": "error-type", invalid json }`
		body := io.NopCloser(bytes.NewReader([]byte(invalidJSON)))

		// Call the function with a prefix
		err := handleErrorResponse(body, "test-prefix")

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "test-prefix")
		assert.Contains(t, err.Error(), "failed to decode error response")

		// Verify it's not an ErrorResponse
		var errorResponse ErrorResponse
		assert.False(t, errors.As(err, &errorResponse))
	})

	t.Run("different prefix", func(t *testing.T) {
		// Create a valid JSON response
		jsonBody := `{
			"type": "error-type",
			"title": "Error Title",
			"status": 400,
			"detail": "Error details",
			"instance": "error-instance"
		}`
		body := io.NopCloser(bytes.NewReader([]byte(jsonBody)))

		// Call the function with a different prefix
		err := handleErrorResponse(body, "another-prefix")

		// Verify the error contains the prefix
		require.Error(t, err)
		assert.Contains(t, err.Error(), "another-prefix")
		assert.NotContains(t, err.Error(), "test-prefix")
	})

	t.Run("empty body", func(t *testing.T) {
		// Create an empty body
		body := io.NopCloser(strings.NewReader(""))

		// Call the function
		err := handleErrorResponse(body, "empty-body")

		// Verify the error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "empty-body")
		assert.Contains(t, err.Error(), "failed to decode error response")
	})
}
