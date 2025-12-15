package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitesinbyte/ferret/pkg/external/groupme"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type sendRequest struct {
	Text string `json:"text"`
}

func main() {
	// Load .env if present (local-only). Non-fatal to keep parity with other cmds.
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}
	// Also try repo root (backend/.env) when running from backend/go.
	_ = godotenv.Load("../.env")

	token := strings.TrimSpace(os.Getenv("GROUPME_WEBHOOK_TOKEN"))
	if token == "" {
		log.Fatal("GROUPME_WEBHOOK_TOKEN is required")
	}
	botID := strings.TrimSpace(os.Getenv("GROUPME_BOT_ID"))
	baseURL := strings.TrimSpace(os.Getenv("GROUPME_BASE_URL"))
	port := strings.TrimSpace(os.Getenv("GROUPME_PORT"))
	if port == "" {
		port = "8081"
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	client := groupme.New(groupme.Config{BaseURL: baseURL})

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Webhook endpoint for GroupMe bot callbacks.
	r.POST("/webhooks/groupme", func(c *gin.Context) {
		ev, err := groupme.ParseAndValidateWebhook(c.Request, groupme.WebhookConfig{Token: token})
		if err != nil {
			switch {
			case err == groupme.ErrUnauthorized:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			}
			return
		}

		log.Printf("groupme webhook: group_id=%s message_id=%s sender_type=%s sender_id=%s name=%q text=%q",
			ev.GroupID, ev.MessageID, ev.SenderType, ev.SenderID, ev.Name, ev.Text,
		)

		// Ignore bot messages (prevents reply loops).
		if ev.SenderType == "bot" || ev.System {
			c.JSON(http.StatusOK, gin.H{"status": "ignored"})
			return
		}

		// Simple demo: reply to !ping with pong (requires GROUPME_BOT_ID).
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(ev.Text)), "!ping") {
			if botID == "" {
				c.JSON(http.StatusNotImplemented, gin.H{"error": "GROUPME_BOT_ID not set"})
				return
			}
			ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
			defer cancel()
			if err := client.PostBotMessage(ctx, botID, "pong"); err != nil {
				log.Printf("groupme post error: %v", err)
				c.JSON(http.StatusBadGateway, gin.H{"error": "failed to post"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Manual outbound endpoint (protected by the same shared secret).
	r.POST("/groupme/send", func(c *gin.Context) {
		if !isAuthorized(c.Request, token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if botID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "GROUPME_BOT_ID not set"})
			return
		}
		var req sendRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		text := strings.TrimSpace(req.Text)
		if text == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing text"})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		if err := client.PostBotMessage(ctx, botID, text); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "sent"})
	})

	addr := ":" + port
	log.Printf("groupmebot listening on %s", addr)
	log.Printf("webhook endpoint: POST http://localhost%s/webhooks/groupme", addr)
	log.Printf("send endpoint:    POST http://localhost%s/groupme/send", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func isAuthorized(r *http.Request, token string) bool {
	exp := strings.TrimSpace(token)
	if exp == "" {
		return false
	}
	if got := strings.TrimSpace(r.URL.Query().Get("token")); got != "" {
		return got == exp
	}
	if got := strings.TrimSpace(r.Header.Get("X-Webhook-Token")); got != "" {
		return got == exp
	}
	if got := strings.TrimSpace(r.Header.Get("Authorization")); got != "" {
		lower := strings.ToLower(got)
		if strings.HasPrefix(lower, "bearer ") {
			return strings.TrimSpace(got[len("Bearer "):]) == exp
		}
	}
	return false
}
