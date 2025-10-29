package rss_test

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bitesinbyte/ferret/pkg/_tests/rss"
	"github.com/joho/godotenv"
)

// TestRealRSSFeed tests fetching and parsing a real RSS feed from a URL
func TestRealRSSFeed(t *testing.T) {
	// Load environment variables from .env file
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Get the feed URL from environment variables
	feedURL := os.Getenv("RSS_FEED_URL")
	if feedURL == "" {
		t.Fatal("RSS_FEED_URL environment variable not set in .env file")
	}
	t.Logf("Using RSS feed URL from .env: %s", feedURL)

	t.Logf("Fetching feed from: %s", feedURL)

	// Create HTTP client with timeout and custom user agent
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set headers to mimic a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml; q=0.9, */*; q=0.8")

	// Fetch the RSS feed
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to fetch RSS feed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Unexpected status code: %d, Response: %s", resp.StatusCode, string(body))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Print the first 1000 characters for debugging
	t.Logf("\n=== RAW XML RESPONSE (first 1000 chars) ===\n%.1000s...\n===========================================\n", string(body))

	// Parse the RSS feed
	var feed rss.RSS
	if err := xml.Unmarshal(body, &feed); err != nil {
		t.Fatalf("Failed to parse RSS feed: %v\nResponse: %s", err, string(body)[:500])
	}

	// Print the parsed RSS structure
	t.Logf("\n=== PARSED RSS STRUCTURE ===")
	t.Logf("Channel Title: %s", feed.Channel.Title)
	channelLink := feed.Channel.GetLink()
	t.Logf("Channel Link: %s", channelLink)
	t.Logf("Channel Description: %s", feed.Channel.Description)
	t.Logf("Last Build Date: %s", feed.Channel.LastBuildDate)
	t.Logf("Generator: %s", feed.Channel.Generator)
	if feed.Channel.AtomLink.Href != "" {
		t.Logf("Atom Link: %s (rel: %s)", feed.Channel.AtomLink.Href, feed.Channel.AtomLink.Rel)
	}
	t.Logf("Number of Items: %d", len(feed.Channel.Items))

	// Print first 3 items (or all if less than 3)
	maxItems := 3
	if len(feed.Channel.Items) < maxItems {
		maxItems = len(feed.Channel.Items)
	}

	for i := 0; i < maxItems; i++ {
		item := feed.Channel.Items[i]
		t.Logf("\n=== Item %d ===", i+1)
		t.Logf("Title: %s", strings.TrimSpace(item.Title))
		t.Logf("Link: %s", item.Link)
		
		// Clean and truncate description
		desc := strings.TrimSpace(item.Description)
		if len(desc) > 150 {
			desc = fmt.Sprintf("%.150s...", desc)
		}
		t.Logf("Description: %s", desc)
		
		t.Logf("Published: %s", item.PubDate)
		if item.Creator != "" {
			t.Logf("Author: %s", item.Creator)
		}
		if len(item.Categories) > 0 {
			t.Logf("Categories: %v", item.Categories)
		}
		if item.Enclosure.URL != "" {
			t.Logf("Enclosure: %s (%s, %d bytes)", 
				item.Enclosure.URL, 
				item.Enclosure.Type, 
				item.Enclosure.Length)
		}
	}

	// Basic validation
	if feed.Channel.Title == "" {
		t.Error("Expected non-empty channel title")
	}

	// Use the existing channelLink variable
	if channelLink == "" {
		t.Error("Expected non-empty channel link")
	}

	if len(feed.Channel.Items) == 0 {
		t.Error("Expected at least one item in the feed")
	} else {
		// Validate first item
		firstItem := feed.Channel.Items[0]
		if firstItem.Title == "" {
			t.Error("Expected first item to have a title")
		}
		if firstItem.Link == "" {
			t.Error("Expected first item to have a link")
		} else {
			t.Logf("First item link: %s", firstItem.Link)
		}
	}

	// Basic validation
	if feed.Channel.Title == "" {
		t.Error("Expected non-empty channel title")
	}

	if len(feed.Channel.Items) == 0 {
		t.Error("Expected at least one item in the feed")
	}
}
