package glif_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	glif "github.com/bitesinbyte/ferret/pkg/api/glif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunWorkflow_PositionalInputs(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()

		var payload map[string]any
		require.NoError(t, json.Unmarshal(body, &payload))
		assert.Equal(t, "workflow-123", payload["id"])
		assert.ElementsMatch(t, []any{"prompt"}, payload["inputs"].([]any))

		_, _ = w.Write([]byte(`{"id":"workflow-123","output":"https://example.com/result.png"}`))
	}))
	t.Cleanup(srv.Close)

	client := glif.NewClient("test-token", glif.WithSimpleAPIBaseURL(srv.URL))
	resp, err := client.RunWorkflow(context.Background(), glif.RunWorkflowRequest{
		WorkflowID: "workflow-123",
		Inputs:     []string{"prompt"},
	})
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/result.png", resp.Output)
}

func TestRunWorkflow_StrictModeAndPathID(t *testing.T) {
	t.Parallel()

	var seenPath, seenQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		seenQuery = r.URL.RawQuery
		assert.Equal(t, "", r.Header.Get("Content-Type"))
		_, _ = w.Write([]byte(`{"id":"workflow-123","output":"ok"}`))
	}))
	t.Cleanup(srv.Close)

	client := glif.NewClient("token", glif.WithSimpleAPIBaseURL(srv.URL))
	resp, err := client.RunWorkflow(context.Background(), glif.RunWorkflowRequest{
		WorkflowID: "workflow-123",
		UsePathID:  true,
		Strict:     true,
	})
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Output)
	assert.Equal(t, "/workflow-123", seenPath)
	assert.Equal(t, "strict=1", seenQuery)
}

func TestRunWorkflow_ReturnsWorkflowError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"id":"workflow-123","error":"missing input"}`))
	}))
	t.Cleanup(srv.Close)

	client := glif.NewClient("token", glif.WithSimpleAPIBaseURL(srv.URL))
	_, err := client.RunWorkflow(context.Background(), glif.RunWorkflowRequest{
		WorkflowID:  "workflow-123",
		NamedInputs: map[string]string{"prompt": "hello"},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, glif.ErrWorkflowFailed)
}

func TestRunWorkflow_ConflictingInputs(t *testing.T) {
	t.Parallel()

	client := glif.NewClient("token")
	_, err := client.RunWorkflow(context.Background(), glif.RunWorkflowRequest{
		WorkflowID:  "workflow-123",
		Inputs:      []string{"a"},
		NamedInputs: map[string]string{"prompt": "b"},
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, glif.ErrConflictingInputModes)
}

func TestRunWorkflow_HTTPError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"invalid token"}`))
	}))
	t.Cleanup(srv.Close)

	client := glif.NewClient("token", glif.WithSimpleAPIBaseURL(srv.URL))
	_, err := client.RunWorkflow(context.Background(), glif.RunWorkflowRequest{
		WorkflowID: "workflow-123",
	})
	require.Error(t, err)
	var apiErr *glif.APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
}

func TestGetWorkflow(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "id=workflow-123", r.URL.RawQuery)
		_, _ = w.Write([]byte(`[{"id":"workflow-123","name":"demo","description":"","output":"", "outputType":"TEXT","user":{"id":"1","name":"Jamie","username":"jamie","image":""},"data":{"nodes":[{"name":"text1","type":"TextInputBlock","params":{"label":"prompt"}}]}}]`))
	}))
	t.Cleanup(srv.Close)

	client := glif.NewClient("token", glif.WithAPIBaseURL(srv.URL))
	wf, err := client.GetWorkflow(context.Background(), "workflow-123")
	require.NoError(t, err)
	assert.Equal(t, "demo", wf.Name)
	assert.Len(t, wf.Data.Nodes, 1)
}

func TestGetWorkflow_NotFound(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[]`))
	}))
	t.Cleanup(srv.Close)

	client := glif.NewClient("token", glif.WithAPIBaseURL(srv.URL))
	_, err := client.GetWorkflow(context.Background(), "missing")
	require.Error(t, err)
	assert.ErrorIs(t, err, glif.ErrWorkflowNotFound)
}
