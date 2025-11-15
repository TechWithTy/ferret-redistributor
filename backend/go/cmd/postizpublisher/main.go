package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/bitesinbyte/ferret/pkg/api/postiz"
	"github.com/bitesinbyte/ferret/pkg/calendar"
	"github.com/bitesinbyte/ferret/pkg/engine/telemetry"
)

func main() {
	var (
		inputPath   = flag.String("input", "go/due_posts.json", "Path to scheduler output (JSON array)")
		dsn         = flag.String("database", os.Getenv("DATABASE_URL"), "Postgres DSN or set DATABASE_URL")
		mapPath     = flag.String("integration-map", "", "Optional path to JSON mapping {\"platform\":\"integrationId\"}")
		skipRewrite = flag.Bool("keep-input", false, "Do not rewrite input file with remaining posts")
	)
	flag.Parse()

	apiKey := strings.TrimSpace(os.Getenv("POSTIZ_API_KEY"))
	if apiKey == "" {
		log.Println("[postizpublisher] POSTIZ_API_KEY not set; skipping.")
		return
	}
	if *dsn == "" {
		log.Fatal("missing database DSN (set --database or DATABASE_URL)")
	}

	rows, err := loadRows(*inputPath)
	if err != nil {
		log.Fatalf("load rows: %v", err)
	}
	if len(rows) == 0 {
		log.Println("[postizpublisher] no rows to process")
		return
	}

	integrationMap, err := loadIntegrationMap(*mapPath)
	if err != nil {
		log.Fatalf("load integration map: %v", err)
	}

	opts := []postiz.Option{postiz.WithAPIKey(apiKey)}
	if base := strings.TrimSpace(os.Getenv("POSTIZ_BASE_URL")); base != "" {
		opts = append(opts, postiz.WithBaseURL(base))
	}
	client := postiz.NewClient(opts...)

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	ctx := telemetry.InitFromEnv(context.Background())
	now := time.Now().UTC()

	cache := make(map[string]postiz.MediaDto)
	var remaining []calendar.ScheduledPostRow
	var processed int
	var successes int

	for _, row := range rows {
		meta := decodeMetadata(row.Metadata)
		integrationID, reason := resolveIntegration(row, meta, integrationMap)
		if integrationID == "" {
			// No integration configured; leave for legacy poster.
			remaining = append(remaining, row)
			if reason != "" {
				log.Printf("[postizpublisher] skip %s: %s", row.ID, reason)
			}
			continue
		}

		ok := handleRow(ctx, client, db, row, integrationID, meta, cache)
		processed++
		if ok {
			successes++
		}
	}

	if !*skipRewrite {
		if err := writeRows(*inputPath, remaining); err != nil {
			log.Printf("[postizpublisher] rewrite %s failed: %v", *inputPath, err)
		}
	}

	log.Printf("[postizpublisher] finished. total=%d processed=%d success=%d skipped=%d in %s",
		len(rows), processed, successes, len(remaining), time.Since(now).Truncate(time.Millisecond))
}

func handleRow(ctx context.Context, client *postiz.Client, db *sql.DB, row calendar.ScheduledPostRow, integrationID string, meta map[string]any, cache map[string]postiz.MediaDto) bool {
	content := composeContent(row, meta)
	if content == "" {
		markFailed(ctx, db, row.ID, "empty content")
		return false
	}

	mediaInputs, err := gatherMedia(ctx, client, meta, cache)
	if err != nil {
		markFailed(ctx, db, row.ID, fmt.Sprintf("media upload failed: %v", err))
		return false
	}

	postType, schedule := determinePostType(row.ScheduledAt)
	req := postiz.CreateUpdatePostRequest{
		Type: postType,
		Date: schedule.Format(time.RFC3339),
		Posts: []postiz.PostInput{
			{
				Integration: postiz.IntegrationInput{ID: integrationID},
				Value: []postiz.PostContent{
					{
						Content: content,
						Image:   mediaInputs,
					},
				},
			},
		},
	}

	res, err := client.Posts.CreateOrUpdate(ctx, req)
	if err != nil {
		markFailed(ctx, db, row.ID, err.Error())
		return false
	}

	var postID string
	if len(res) > 0 {
		postID = res[0].PostID
	}

	updates := map[string]any{
		"postiz": map[string]any{
			"post_id":      postID,
			"integration":  integrationID,
			"state":        "queue",
			"submitted_at": time.Now().UTC().Format(time.RFC3339),
		},
	}
	if len(mediaInputs) > 0 {
		urls := make([]string, 0, len(mediaInputs))
		for _, m := range mediaInputs {
			if m.Path != "" {
				urls = append(urls, m.Path)
			}
		}
		if len(urls) > 0 {
			updates["postiz"].(map[string]any)["media_paths"] = urls
		}
	}

	var publishedAt *time.Time
	if postType == "now" {
		ts := time.Now().UTC()
		publishedAt = &ts
	}
	merged := mergeMetadata(row.Metadata, updates)
	if err := calendar.UpdatePostStatus(ctx, db, row.ID, calendar.StatusPublished, stringPtr(postID), publishedAt, merged); err != nil {
		log.Printf("[postizpublisher] update status %s failed: %v", row.ID, err)
		return false
	}
	log.Printf("[postizpublisher] queued %s (%s) via Postiz", row.ID, integrationID)
	return true
}

func loadRows(path string) ([]calendar.ScheduledPostRow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rows []calendar.ScheduledPostRow
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func writeRows(path string, rows []calendar.ScheduledPostRow) error {
	data, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func loadIntegrationMap(path string) (map[string]string, error) {
	var data []byte
	var err error
	if path != "" {
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, err
		}
	} else if raw := strings.TrimSpace(os.Getenv("POSTIZ_INTEGRATIONS")); raw != "" {
		data = []byte(raw)
	}
	if len(data) == 0 {
		return map[string]string{}, nil
	}
	rawMap := map[string]string{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rawMap))
	for k, v := range rawMap {
		key := strings.ToLower(strings.TrimSpace(k))
		if key == "" {
			continue
		}
		out[key] = strings.TrimSpace(v)
	}
	return out, nil
}

func resolveIntegration(row calendar.ScheduledPostRow, meta map[string]any, integrations map[string]string) (string, string) {
	if meta != nil {
		if id := firstString(meta, "postiz_integration", "postizIntegration", "integration_id", "postizIntegrationId"); id != "" {
			return id, ""
		}
	}
	platform := strings.ToLower(string(row.Platform))
	if id := integrations[platform]; id != "" {
		return id, ""
	}
	if id := integrations["*"]; id != "" {
		return id, ""
	}
	return "", fmt.Sprintf("no Postiz integration configured for %s", platform)
}

func gatherMedia(ctx context.Context, client *postiz.Client, meta map[string]any, cache map[string]postiz.MediaDto) ([]postiz.MediaDto, error) {
	urls := mediaURLs(meta)
	if len(urls) == 0 {
		return nil, nil
	}
	out := make([]postiz.MediaDto, 0, len(urls))
	for _, u := range urls {
		if cached, ok := cache[u]; ok {
			out = append(out, cached)
			continue
		}
		resp, err := client.Upload.UploadFromURL(ctx, postiz.UploadFromURLRequest{URL: u})
		if err != nil {
			return nil, err
		}
		dto := postiz.MediaDto{ID: resp.ID, Path: resp.Path}
		cache[u] = dto
		out = append(out, dto)
	}
	return out, nil
}

func mediaURLs(meta map[string]any) []string {
	if meta == nil {
		return nil
	}
	var collected []string
	for _, key := range []string{"media_urls", "image_urls", "images", "media"} {
		if v, ok := meta[key]; ok {
			collected = append(collected, expandStringValues(v)...)
		}
	}
	for _, key := range []string{"image_url", "media_url"} {
		if v := firstString(meta, key); v != "" {
			collected = append(collected, v)
		}
	}
	return uniqueStrings(collected)
}

func expandStringValues(v any) []string {
	switch val := v.(type) {
	case []any:
		out := make([]string, 0, len(val))
		for _, item := range val {
			if s := fmt.Sprint(item); strings.TrimSpace(s) != "" {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out
	case []string:
		out := make([]string, 0, len(val))
		for _, s := range val {
			if strings.TrimSpace(s) != "" {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out
	default:
		if s := strings.TrimSpace(fmt.Sprint(val)); s != "" {
			return []string{s}
		}
	}
	return nil
}

func composeContent(row calendar.ScheduledPostRow, meta map[string]any) string {
	if meta != nil {
		if v := firstString(meta, "postiz_content", "content_override"); v != "" {
			return v
		}
	}

	var parts []string
	if caption := strings.TrimSpace(nullString(row.Caption)); caption != "" {
		parts = append(parts, caption)
	} else if meta != nil {
		if caption := firstString(meta, "caption"); caption != "" {
			parts = append(parts, caption)
		}
	}
	if tags := strings.TrimSpace(nullString(row.Hashtags)); tags != "" {
		parts = append(parts, tags)
	}
	if link := strings.TrimSpace(nullString(row.ContentURL)); link != "" {
		parts = append(parts, link)
	}
	if len(parts) == 0 {
		if row.ContentTitle.Valid {
			parts = append(parts, row.ContentTitle.String)
		} else if strings.TrimSpace(row.CampaignName) != "" {
			parts = append(parts, row.CampaignName)
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n\n"))
}

func determinePostType(ts time.Time) (string, time.Time) {
	if ts.IsZero() {
		return "now", time.Now().UTC()
	}
	when := ts.UTC()
	if when.Before(time.Now().UTC().Add(-1 * time.Minute)) {
		return "now", time.Now().UTC()
	}
	return "schedule", when
}

func decodeMetadata(raw json.RawMessage) map[string]any {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil
	}
	return out
}

func mergeMetadata(original json.RawMessage, updates map[string]any) json.RawMessage {
	if updates == nil {
		return original
	}
	var merged map[string]any
	if len(original) > 0 && string(original) != "null" {
		if err := json.Unmarshal(original, &merged); err != nil {
			merged = map[string]any{"_raw_metadata": string(original)}
		}
	}
	if merged == nil {
		merged = make(map[string]any, len(updates))
	}
	for k, v := range updates {
		merged[k] = v
	}
	data, err := json.Marshal(merged)
	if err != nil {
		return original
	}
	return data
}

func markFailed(ctx context.Context, db *sql.DB, id, msg string) {
	meta := map[string]any{"error": msg, "postiz": map[string]any{"state": "failed"}}
	if err := calendar.UpdatePostStatus(ctx, db, id, calendar.StatusFailed, nil, nil, marshalMap(meta)); err != nil {
		log.Printf("[postizpublisher] update failed status %s: %v (reason: %s)", id, err, msg)
	} else {
		log.Printf("[postizpublisher] marked %s failed: %s", id, msg)
	}
}

func marshalMap(m map[string]any) json.RawMessage {
	if m == nil {
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return b
}

func firstString(meta map[string]any, keys ...string) string {
	if meta == nil {
		return ""
	}
	for _, key := range keys {
		if val, ok := meta[key]; ok {
			if s, err := toString(val); err == nil && s != "" {
				return s
			}
		}
	}
	return ""
}

func toString(v any) (string, error) {
	switch val := v.(type) {
	case string:
		return strings.TrimSpace(val), nil
	case fmt.Stringer:
		return strings.TrimSpace(val.String()), nil
	case float64:
		return strings.TrimSpace(fmt.Sprintf("%.0f", val)), nil
	case json.Number:
		return strings.TrimSpace(val.String()), nil
	case nil:
		return "", nil
	default:
		return strings.TrimSpace(fmt.Sprint(val)), nil
	}
}

func uniqueStrings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func stringPtr(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}
