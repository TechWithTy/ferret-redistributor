package fal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallAndPollUntilCompletion(t *testing.T) {
	t.Parallel()

	var mu sync.Mutex
	statusCalls := 0

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/fal/test":
			assert.Equal(t, http.MethodPost, r.Method)
			writeJSON(t, w, Response{
				RequestID:   "req-1",
				StatusURL:   serverURLPath(server, "/status/req-1"),
				ResponseURL: serverURLPath(server, "/response/req-1"),
				CancelURL:   serverURLPath(server, "/cancel/req-1"),
			})
		case "/status/req-1":
			mu.Lock()
			defer mu.Unlock()
			statusCalls++
			if statusCalls == 1 {
				writeJSON(t, w, Response{Status: "IN_PROGRESS"})
				return
			}
			writeJSON(t, w, Response{Status: "COMPLETED"})
		case "/response/req-1":
			writeJSON(t, w, Response{
				Status:    "COMPLETED",
				RequestID: "req-1",
				Images: []Image{
					{URL: "https://fal.media/example.png"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := NewClient("secret",
		WithQueueBaseURL(server.URL),
		WithPollInterval(10*time.Millisecond),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	call, err := client.Call(ctx, "fal/test", RequestOptions{
		Payload: map[string]string{"prompt": "hello"},
	})
	require.NoError(t, err)

	resp, err := client.PollUntilCompletion(ctx, call)
	require.NoError(t, err)
	assert.Equal(t, "req-1", resp.RequestID)
	require.Len(t, resp.Images, 1)
	assert.Equal(t, "https://fal.media/example.png", resp.Images[0].URL)
}

func TestPollWithProgress(t *testing.T) {
	t.Parallel()

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/fal/test":
			writeJSON(t, w, Response{
				RequestID:   "req-2",
				StatusURL:   serverURLPath(server, "/status/req-2"),
				ResponseURL: serverURLPath(server, "/response/req-2"),
			})
		case "/status/req-2":
			writeJSON(t, w, Response{Status: "COMPLETED", QueuePosition: 0})
		case "/response/req-2":
			writeJSON(t, w, Response{Status: "COMPLETED"})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := NewClient("secret",
		WithQueueBaseURL(server.URL),
		WithPollInterval(5*time.Millisecond),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	call, err := client.Call(ctx, "fal/test", RequestOptions{})
	require.NoError(t, err)

	progressCh := make(chan *Response, 1)
	doneCh := make(chan struct{})
	go func() {
		for range progressCh {
		}
		close(doneCh)
	}()

	resp, err := client.PollWithProgress(ctx, call, progressCh)
	close(progressCh)
	<-doneCh

	require.NoError(t, err)
	assert.Equal(t, "COMPLETED", resp.Status)
}

func TestCancel(t *testing.T) {
	t.Parallel()

	cancelCalled := false
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/fal/test":
			writeJSON(t, w, Response{
				RequestID: "req-3",
				CancelURL: serverURLPath(server, "/cancel/req-3"),
			})
		case "/cancel/req-3":
			cancelCalled = true
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := NewClient("secret", WithQueueBaseURL(server.URL))
	ctx := context.Background()

	call, err := client.Call(ctx, "fal/test", RequestOptions{})
	require.NoError(t, err)

	require.NoError(t, call.Cancel(ctx))
	assert.True(t, cancelCalled)
}

func TestPollUntilCompletionFailure(t *testing.T) {
	t.Parallel()

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/fal/test":
			writeJSON(t, w, Response{
				RequestID: "req-4",
				StatusURL: serverURLPath(server, "/status/req-4"),
			})
		case "/status/req-4":
			writeJSON(t, w, Response{Status: "FAILED"})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := NewClient("secret", WithQueueBaseURL(server.URL), WithPollInterval(5*time.Millisecond))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	call, err := client.Call(ctx, "fal/test", RequestOptions{})
	require.NoError(t, err)

	_, err = client.PollUntilCompletion(ctx, call)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrCallIncomplete)
}

func writeJSON(t *testing.T, w http.ResponseWriter, payload any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(w).Encode(payload))
}

func serverURLPath(server *httptest.Server, path string) string {
	return server.URL + path
}
