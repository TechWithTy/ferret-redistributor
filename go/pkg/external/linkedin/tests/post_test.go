package linkedin_test

import (
	"os"
	"testing"
	"time"

	"github.com/bitesinbyte/ferret/pkg/config"
	"github.com/bitesinbyte/ferret/pkg/external/linkedin"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	// Load environment variables
	accessToken := os.Getenv("LINKEDIN_ACCESS_TOKEN")
	refreshToken := os.Getenv("LINKEDIN_REFRESH_TOKEN")
	clientID := os.Getenv("LINKEDIN_CLIENT_ID")
	clientSecret := os.Getenv("LINKEDIN_CLIENT_SECRET")

	// Skip test if required credentials are not set
	if accessToken == "" || refreshToken == "" || clientID == "" || clientSecret == "" {
		t.Skip("Skipping test: LinkedIn credentials not fully configured")
	}

	// For now, we'll just use the access token directly
	// since the LinkedIn package handles the HTTP client internally
	os.Setenv("LINKEDIN_ACCESS_TOKEN", accessToken)

	// Initialize config with required fields
	cfg := config.Config{
		BaseUrl: "https://example.com",
	}

	// Initialize LinkedIn client
	client := linkedin.Linkedin{}

	// Test post content
	postContent := "This is a test post created by Ferret at " + time.Now().Format(time.RFC3339)

	// Create post
	post := linkedin.Post{
		Title:       "Test Post",
		Description: postContent,
		Link:        "https://example.com",
		HashTags:    []string{"#test", "#golang"},
	}

	// Test PostWithID which returns the post URN
	t.Run("CreatePostWithID", func(t *testing.T) {
		postID, err := client.PostWithID(cfg, post)
		require.NoError(t, err, "Failed to create post with ID")
		t.Logf("Created post with ID: %s", postID)
	})

	// Test regular Post method
	t.Run("CreatePost", func(t *testing.T) {
		err := client.Post(cfg, post)
		require.NoError(t, err, "Failed to create post")
		t.Log("Successfully created post")
	})
}
