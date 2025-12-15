package notion_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitesinbyte/ferret/pkg/external/notion"
	"github.com/stretchr/testify/require"
)

func TestClient_QueryFirstByTitle_UsesDataSourcesEndpointAndVersionHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer ntn_test", r.Header.Get("Authorization"))
		require.Equal(t, "2025-09-03", r.Header.Get("Notion-Version"))
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/data_sources/ds123/query", r.URL.Path)

		var req map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		// basic shape check
		filter := req["filter"].(map[string]any)
		require.Equal(t, "Group ID", filter["property"])
		title := filter["title"].(map[string]any)
		require.Equal(t, "g1", title["equals"])

		_, _ = w.Write([]byte(`{"results":[{"id":"p1","url":"https://notion.so/p1"}]}`))
	}))
	t.Cleanup(srv.Close)

	c, err := notion.New(notion.Config{APIKey: "ntn_test", BaseURL: srv.URL + "/v1"})
	require.NoError(t, err)

	p, err := c.QueryFirstByTitle(context.Background(), "ds123", "Group ID", "g1")
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, "p1", p.ID)
}

func TestClient_CreatePageInDataSource_ParentUsesDataSourceID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/pages", r.URL.Path)

		var req map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		parent := req["parent"].(map[string]any)
		require.Equal(t, "data_source_id", parent["type"])
		require.Equal(t, "ds456", parent["data_source_id"])

		_, _ = w.Write([]byte(`{"id":"p2","url":"https://notion.so/p2"}`))
	}))
	t.Cleanup(srv.Close)

	c, err := notion.New(notion.Config{APIKey: "ntn_test", BaseURL: srv.URL + "/v1"})
	require.NoError(t, err)

	created, err := c.CreatePageInDataSource(context.Background(), "ds456", map[string]any{
		"Group ID": notion.Title("g1"),
	})
	require.NoError(t, err)
	require.Equal(t, "p2", created.ID)
}

func TestClient_UpdatePageProperties_Patch(t *testing.T) {
	var gotAuth, gotVer string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotVer = r.Header.Get("Notion-Version")
		require.Equal(t, http.MethodPatch, r.Method)
		require.True(t, strings.HasSuffix(r.URL.Path, "/v1/pages/p3"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(srv.Close)

	c, err := notion.New(notion.Config{APIKey: "ntn_test", BaseURL: srv.URL + "/v1"})
	require.NoError(t, err)

	require.NoError(t, c.UpdatePageProperties(context.Background(), "p3", map[string]any{
		"Group Name": notion.RichText("Name"),
	}))
	require.Equal(t, "Bearer ntn_test", gotAuth)
	require.Equal(t, "2025-09-03", gotVer)
}
