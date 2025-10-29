package rss_test

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bitesinbyte/ferret/pkg/_tests/rss"
)

// createTestServer creates a test HTTP server that serves a sample RSS feed
func createTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sampleRSS := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test Feed</title>
    <link>https://example.com</link>
    <description>This is a test RSS feed</description>
    <item>
      <title>Test Item 1</title>
      <link>https://example.com/item1</link>
      <description>This is a test item</description>
      <pubDate>%s</pubDate>
    </item>
  </channel>
</rss>`
		// Use current time for the pubDate
		sampleRSS = fmt.Sprintf(sampleRSS, time.Now().Format(time.RFC1123Z))
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sampleRSS))
	}))

	t.Cleanup(func() { server.Close() })
	return server
}

// TestFetchRSSFeed demonstrates how to fetch and parse an RSS feed
func TestFetchRSSFeed(t *testing.T) {
	t.Run("should parse RSS feed from test server", func(t *testing.T) {
		testServer := createTestServer(t)

		// Fetch the RSS feed
		resp, err := http.Get(testServer.URL)
		if err != nil {
			t.Fatalf("Failed to fetch RSS feed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		// Parse the RSS feed
		var feed rss.RSS
		err = xml.Unmarshal(body, &feed)
		if err != nil {
			t.Fatalf("Failed to parse RSS feed: %v", err)
		}

		// Verify the feed was parsed correctly
		if feed.Channel.Title != "Test Feed" {
			t.Errorf("Expected title 'Test Feed', got '%s'", feed.Channel.Title)
		}

		if len(feed.Channel.Items) == 0 {
			t.Error("Expected at least one item in the feed")
		} else {
			item := feed.Channel.Items[0]
			if item.Title != "Test Item 1" {
				t.Errorf("Expected item title 'Test Item 1', got '%s'", item.Title)
			}
		}
	})

	t.Run("should use SITEMAP_URL from environment variable", func(t *testing.T) {
		// Set up test server
		testServer := createTestServer(t)

		// Set SITEMAP_URL environment variable for this test
		t.Setenv("SITEMAP_URL", testServer.URL)

		// In your actual code, you would use os.Getenv("SITEMAP_URL") to get the URL
		sitemapURL := os.Getenv("SITEMAP_URL")
		if sitemapURL == "" {
			t.Fatal("SITEMAP_URL environment variable not set")
		}

		// Fetch the RSS feed using the URL from SITEMAP_URL
		resp, err := http.Get(sitemapURL)
		if err != nil {
			t.Fatalf("Failed to fetch RSS feed from SITEMAP_URL: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Unexpected status code from SITEMAP_URL: %d", resp.StatusCode)
		}

		// Read and parse the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		// Print the raw XML for debugging
		t.Logf("\n=== RAW XML RESPONSE ===\n%s\n========================\n", string(body))

		var feed rss.RSS
		if err := xml.Unmarshal(body, &feed); err != nil {
			t.Fatalf("Failed to parse RSS feed from SITEMAP_URL: %v", err)
		}

		// Print the parsed RSS structure
		t.Logf("\n=== PARSED RSS STRUCTURE ===")
		t.Logf("Channel Title: %s", feed.Channel.Title)
		t.Logf("Channel Link: %s", feed.Channel.Link)
		t.Logf("Channel Description: %s", feed.Channel.Description)
		t.Logf("Number of Items: %d", len(feed.Channel.Items))

		for i, item := range feed.Channel.Items {
			t.Logf("\nItem %d:", i+1)
			t.Logf("  Title: %s", item.Title)
			t.Logf("  Link: %s", item.Link)
			t.Logf("  Description: %s", item.Description)
			t.Logf("  PubDate: %s", item.PubDate)
		}

		// Basic validation
		if feed.Channel.Title != "Test Feed" {
			t.Errorf("Expected title 'Test Feed', got '%s'", feed.Channel.Title)
		}

		if len(feed.Channel.Items) == 0 {
			t.Error("Expected at least one item in the feed from SITEMAP_URL")
		}
	})
}
