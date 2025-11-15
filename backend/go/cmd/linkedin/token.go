package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

// Config holds the OAuth2 configuration
type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

func main() {
	// Load configuration from environment variables or config file
	config := oauth2.Config{
		ClientID:     os.Getenv("LINKEDIN_CLIENT_ID"),
		ClientSecret: os.Getenv("LINKEDIN_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/linkedin/callback",
		Scopes:       []string{"r_liteprofile", "r_emailaddress", "w_member_social"},
		Endpoint:     linkedin.Endpoint,
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		log.Fatal("LINKEDIN_CLIENT_ID and LINKEDIN_CLIENT_SECRET environment variables must be set")
	}

	// Generate the authorization URL
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit this URL to authorize the application:\n\n%s\n\n", url)

	// Start a temporary server to handle the OAuth2 callback
	code := make(chan string)
	http.HandleFunc("/auth/linkedin/callback", func(w http.ResponseWriter, r *http.Request) {
		code <- r.URL.Query().Get("code")
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body><h1>Authentication successful! You can close this window.</h1></body></html>`))
	})

	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Exchange the authorization code for an access token
	authCode := <-code
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Error exchanging code for token: %v", err)
	}

	// Print the token in JSON format
	tokenJSON, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling token: %v", err)
	}

	fmt.Println("\nSuccess! Here's your access token:")
	fmt.Println(string(tokenJSON))

	// Save to .env file
	envContent := fmt.Sprintf(`LINKEDIN_ACCESS_TOKEN=%s
LINKEDIN_CLIENT_ID=%s
LINKEDIN_CLIENT_SECRET=%s
`, token.AccessToken, config.ClientID, config.ClientSecret)

	if err := os.WriteFile(".env", []byte(envContent), 0644); err != nil {
		log.Printf("Warning: Could not write to .env file: %v", err)
	} else {
		fmt.Println("\nToken saved to .env file")
	}
}
