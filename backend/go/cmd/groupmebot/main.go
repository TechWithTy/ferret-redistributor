package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bitesinbyte/ferret/pkg/external/groupme"
	"github.com/bitesinbyte/ferret/pkg/external/notion"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type sendRequest struct {
	Text string `json:"text"`
}

type notionMessageLogger struct {
	enabled bool

	nc *notion.Client

	dsGroups string
	dsBots   string
	dsLogs   string

	botID     string
	botPageID string

	mu                  sync.RWMutex
	groupPageIDByGroupID map[string]string
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

	nlog := newNotionMessageLogger(botID)
	if nlog.enabled {
		log.Printf("notion logging enabled: Bot Message Logs=%s", nlog.dsLogs)
	} else {
		log.Printf("notion logging disabled (set NOTION_API_KEY + NOTION_DATA_SOURCE_ID_BOT_MESSAGE_LOGS to enable)")
	}

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

		// Best-effort Notion logging (non-blocking).
		nlog.logInbound(ev)

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

			// Best-effort Notion logging for outbound reply (non-blocking).
			nlog.logOutbound(botID, ev.GroupID, "pong", "reply_to:"+ev.MessageID)
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

		// Best-effort Notion logging (non-blocking). No group_id available here.
		nlog.logOutbound(botID, "", text, fmt.Sprintf("manual:%d", time.Now().UTC().UnixNano()))

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

func newNotionMessageLogger(botID string) *notionMessageLogger {
	key := strings.TrimSpace(getenvAny("NOTION_API_KEY", "NOTION_TOKEN", "NOTION_KEY"))
	dsLogs := strings.TrimSpace(os.Getenv("NOTION_DATA_SOURCE_ID_BOT_MESSAGE_LOGS"))
	if key == "" || dsLogs == "" {
		return &notionMessageLogger{enabled: false}
	}

	dsGroups := strings.TrimSpace(os.Getenv("NOTION_DATA_SOURCE_ID_GROUPS"))
	dsBots := strings.TrimSpace(os.Getenv("NOTION_DATA_SOURCE_ID_BOTS"))

	nc, err := notion.New(notion.Config{APIKey: key, HTTPTimeout: 30 * time.Second})
	if err != nil {
		log.Printf("notion logger init error: %v", err)
		return &notionMessageLogger{enabled: false}
	}

	l := &notionMessageLogger{
		enabled:              true,
		nc:                   nc,
		dsGroups:             dsGroups,
		dsBots:               dsBots,
		dsLogs:               dsLogs,
		botID:                strings.TrimSpace(botID),
		groupPageIDByGroupID: make(map[string]string, 64),
	}

	// Resolve bot page id once (best-effort).
	if l.botID != "" && l.dsBots != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ref, err := l.nc.QueryFirstByTitle(ctx, l.dsBots, "Bot ID", l.botID)
		if err != nil {
			log.Printf("notion logger: failed resolving bot page: %v", err)
		} else if ref != nil {
			l.botPageID = ref.ID
		}
	}

	return l
}

func (l *notionMessageLogger) logInbound(ev groupme.WebhookEvent) {
	if l == nil || !l.enabled || l.nc == nil {
		return
	}

	// Webhook must stay fast; do the Notion write async.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		logID := messageLogID("Response", l.botID, ev.GroupID, ev.MessageID)
		props := map[string]any{
			"Log ID":       notion.Title(logID),
			"Direction":    notion.Select("Response"),
			"Message ID":   notion.RichText(ev.MessageID),
			"Message Text": notion.RichText(ev.Text),
		}
		if ev.CreatedAt > 0 {
			props["Timestamp"] = notion.DateTime(time.Unix(ev.CreatedAt, 0))
		} else {
			props["Timestamp"] = notion.DateTime(time.Now().UTC())
		}

		if strings.TrimSpace(l.botPageID) != "" {
			props["Bot"] = notion.Relation(l.botPageID)
		}
		if gid := strings.TrimSpace(ev.GroupID); gid != "" {
			if pid := l.resolveGroupPageID(ctx, gid); pid != "" {
				props["Group"] = notion.Relation(pid)
			}
		}

		if _, err := l.nc.UpsertByTitle(ctx, l.dsLogs, "Log ID", logID, props); err != nil {
			log.Printf("notion logger: inbound upsert failed: %v", err)
		}
	}()
}

func (l *notionMessageLogger) logOutbound(botID, groupID, text, idempotencyKey string) {
	if l == nil || !l.enabled || l.nc == nil {
		return
	}
	if strings.TrimSpace(text) == "" {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		logID := messageLogID("Outbound", botID, groupID, idempotencyKey)
		props := map[string]any{
			"Log ID":       notion.Title(logID),
			"Direction":    notion.Select("Outbound"),
			"Message Text": notion.RichText(text),
			"Timestamp":    notion.DateTime(time.Now().UTC()),
		}

		if strings.TrimSpace(l.botPageID) != "" {
			props["Bot"] = notion.Relation(l.botPageID)
		}
		if gid := strings.TrimSpace(groupID); gid != "" {
			if pid := l.resolveGroupPageID(ctx, gid); pid != "" {
				props["Group"] = notion.Relation(pid)
			}
		}

		if _, err := l.nc.UpsertByTitle(ctx, l.dsLogs, "Log ID", logID, props); err != nil {
			log.Printf("notion logger: outbound upsert failed: %v", err)
		}
	}()
}

func (l *notionMessageLogger) resolveGroupPageID(ctx context.Context, groupID string) string {
	if strings.TrimSpace(groupID) == "" || strings.TrimSpace(l.dsGroups) == "" {
		return ""
	}
	l.mu.RLock()
	if pid := l.groupPageIDByGroupID[groupID]; pid != "" {
		l.mu.RUnlock()
		return pid
	}
	l.mu.RUnlock()

	ref, err := l.nc.QueryFirstByTitle(ctx, l.dsGroups, "Group ID", groupID)
	if err != nil || ref == nil || strings.TrimSpace(ref.ID) == "" {
		return ""
	}
	l.mu.Lock()
	l.groupPageIDByGroupID[groupID] = ref.ID
	l.mu.Unlock()
	return ref.ID
}

func messageLogID(direction, botID, groupID, messageOrKey string) string {
	// Keep it deterministic and Notion-title-friendly.
	// (Also hash long bits to avoid extremely long titles.)
	d := strings.TrimSpace(direction)
	if d == "" {
		d = "Unknown"
	}
	b := strings.TrimSpace(botID)
	g := strings.TrimSpace(groupID)
	m := strings.TrimSpace(messageOrKey)
	raw := fmt.Sprintf("%s|%s|%s|%s", d, b, g, m)
	h := sha1.Sum([]byte(raw))
	return fmt.Sprintf("groupme:%s:%s", d, hex.EncodeToString(h[:]))
}

func getenvAny(keys ...string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return ""
}
