package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bitesinbyte/ferret/pkg/external/groupme"
	"github.com/bitesinbyte/ferret/pkg/external/notion"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env from current dir and from repo root (backend/.env) when running from backend/go.
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	var (
		apply       = flag.Bool("apply", false, "write changes to Notion (default: dry-run)")
		metrics     = flag.Bool("metrics", false, "also update Group Engagement Metrics from Bot Message Logs")
		metricsOnly = flag.Bool("metrics-only", false, "only update metrics (skip GroupMe groups/bots sync)")
		windowDays  = flag.Int("window-days", 30, "metrics window size in days (current period)")
	)
	flag.Parse()

	// Notion token
	notionKey := getenvAny("NOTION_API_KEY", "NOTION_TOKEN", "NOTION_KEY")
	if notionKey == "" {
		log.Fatal("missing Notion API key (set NOTION_API_KEY or NOTION_TOKEN)")
	}

	// Notion data sources (from your DB stack)
	dsGroups := getenvAny("NOTION_DATA_SOURCE_ID_GROUPS")
	dsBots := getenvAny("NOTION_DATA_SOURCE_ID_BOTS")
	if dsGroups == "" || dsBots == "" {
		log.Fatal("missing Notion data source IDs (set NOTION_DATA_SOURCE_ID_GROUPS and NOTION_DATA_SOURCE_ID_BOTS)")
	}

	// This command can perform many Notion writes; keep a generous overall timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	nc, err := notion.New(notion.Config{APIKey: notionKey, HTTPTimeout: 60 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("mode: %s", ternary(*apply, "APPLY", "DRY-RUN"))

	// Optionally: sync GroupMe groups/bots into Notion.
	if !*metricsOnly {
		// GroupMe token: accept both names (your repo already had GROUP_ME_MASTER_TOKEN).
		gmToken := getenvAny("GROUPME_ACCESS_TOKEN", "GROUP_ME_MASTER_TOKEN")
		if gmToken == "" {
			log.Fatal("missing GroupMe token (set GROUPME_ACCESS_TOKEN or GROUP_ME_MASTER_TOKEN)")
		}
		gmc, err := groupme.NewUserClient(groupme.UserConfig{AccessToken: gmToken})
		if err != nil {
			log.Fatal(err)
		}

		groups, err := gmc.ListGroups(ctx)
		if err != nil {
			log.Fatal(err)
		}
		bots, err := gmc.ListBots(ctx)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("groupme: %d groups, %d bots fetched", len(groups), len(bots))

		// Load existing Notion pages once (avoid per-item query).
		existingGroups, err := nc.QueryPageRefsByTitle(ctx, dsGroups, "Group ID")
		if err != nil {
			log.Fatalf("failed to load Notion Groups index: %v", err)
		}
		existingBots, err := nc.QueryPageRefsByTitle(ctx, dsBots, "Bot ID")
		if err != nil {
			log.Fatalf("failed to load Notion Bots index: %v", err)
		}

		// 1) Upsert groups first, build map group_id -> notion page id
		groupPageIDByGroupID := make(map[string]string, len(existingGroups)+len(groups))
		for gid, ref := range existingGroups {
			groupPageIDByGroupID[gid] = ref.ID
		}

		for _, g := range groups {
			props := map[string]any{
				"Group ID":      notion.Title(g.ID),
				"Members Count": notion.Number(float64(g.MembersCount)),
			}
			if strings.TrimSpace(g.Name) != "" {
				props["Group Name"] = notion.RichText(g.Name)
			}
			if strings.TrimSpace(g.CreatorUserID) != "" {
				props["Owner ID"] = notion.RichText(g.CreatorUserID)
			}
			if g.UpdatedAt > 0 {
				props["Last Activity"] = notion.DateTime(time.Unix(g.UpdatedAt, 0))
			}
			if ref, ok := existingGroups[g.ID]; ok && strings.TrimSpace(ref.ID) != "" {
				if !*apply {
					log.Printf("[dry-run] update Group ID=%s page_id=%s", g.ID, ref.ID)
				} else {
					if err := nc.UpdatePageProperties(ctx, ref.ID, props); err != nil {
						log.Fatalf("groups upsert failed (group_id=%s): %v", g.ID, err)
					}
					log.Printf("[update] Group ID=%s page_id=%s", g.ID, ref.ID)
				}
				continue
			}
			if !*apply {
				log.Printf("[dry-run] create Group ID=%s", g.ID)
				continue
			}
			created, err := nc.CreatePageInDataSource(ctx, dsGroups, props)
			if err != nil {
				log.Fatalf("groups upsert failed (group_id=%s): %v", g.ID, err)
			}
			groupPageIDByGroupID[g.ID] = created.ID
			log.Printf("[create] Group ID=%s page_id=%s", g.ID, created.ID)
		}

		// 2) Upsert bots, link to group relation if possible
		now := time.Now().UTC()
		for _, b := range bots {
			props := map[string]any{
				"Bot ID":      notion.Title(b.BotID),
				"Active":      notion.Checkbox(true),
				"Last Synced": notion.DateTime(now),
			}
			if strings.TrimSpace(b.Name) != "" {
				props["Bot Name"] = notion.RichText(b.Name)
			}
			if isValidHTTPURL(b.AvatarURL) {
				props["Avatar URL"] = notion.URL(b.AvatarURL)
			}
			if isValidHTTPURL(b.CallbackURL) {
				props["Callback URL"] = notion.URL(b.CallbackURL)
			}
			if b.CreatedAt > 0 {
				props["Created At"] = notion.DateTime(time.Unix(b.CreatedAt, 0))
			}
			if gid := strings.TrimSpace(b.GroupID); gid != "" {
				if pid := groupPageIDByGroupID[gid]; pid != "" {
					props["Group Relation"] = notion.Relation(pid)
				}
			}
			if ref, ok := existingBots[b.BotID]; ok && strings.TrimSpace(ref.ID) != "" {
				if !*apply {
					log.Printf("[dry-run] update Bot ID=%s page_id=%s", b.BotID, ref.ID)
				} else {
					if err := nc.UpdatePageProperties(ctx, ref.ID, props); err != nil {
						log.Fatalf("bots upsert failed (bot_id=%s): %v", b.BotID, err)
					}
					log.Printf("[update] Bot ID=%s page_id=%s", b.BotID, ref.ID)
				}
				continue
			}
			if !*apply {
				log.Printf("[dry-run] create Bot ID=%s", b.BotID)
				continue
			}
			created, err := nc.CreatePageInDataSource(ctx, dsBots, props)
			if err != nil {
				log.Fatalf("bots upsert failed (bot_id=%s): %v", b.BotID, err)
			}
			log.Printf("[create] Bot ID=%s page_id=%s", b.BotID, created.ID)
		}
	}

	// Optionally: metrics update from Bot Message Logs -> Group Engagement Metrics.
	if *metrics || *metricsOnly {
		if err := runMetrics(ctx, nc, *apply, dsGroups, *windowDays); err != nil {
			log.Fatalf("metrics update failed: %v", err)
		}
	}

	log.Printf("done")
}

func runMetrics(ctx context.Context, nc *notion.Client, apply bool, dsGroups string, windowDays int) error {
	dsLogs := getenvAny("NOTION_DATA_SOURCE_ID_BOT_MESSAGE_LOGS")
	dsMetrics := getenvAny("NOTION_DATA_SOURCE_ID_GROUP_ENGAGEMENT_METRICS")
	if dsLogs == "" || dsMetrics == "" {
		return fmt.Errorf("missing Notion data source IDs for metrics (set NOTION_DATA_SOURCE_ID_BOT_MESSAGE_LOGS and NOTION_DATA_SOURCE_ID_GROUP_ENGAGEMENT_METRICS)")
	}
	if windowDays <= 0 {
		windowDays = 30
	}

	// Build group page-id -> group-id map for stable metric titles.
	groupsIndex, err := nc.QueryPageRefsByTitle(ctx, dsGroups, "Group ID")
	if err != nil {
		return fmt.Errorf("load groups index: %w", err)
	}
	groupIDByPageID := make(map[string]string, len(groupsIndex))
	for gid, ref := range groupsIndex {
		if strings.TrimSpace(ref.ID) == "" {
			continue
		}
		groupIDByPageID[ref.ID] = gid
	}

	metricsIndex, err := nc.QueryPageRefsByTitle(ctx, dsMetrics, "Name")
	if err != nil {
		return fmt.Errorf("load metrics index: %w", err)
	}

	now := time.Now().UTC()
	curStart := now.AddDate(0, 0, -windowDays)
	prevStart := curStart.AddDate(0, 0, -windowDays)

	type agg struct {
		responses   int
		prevOutbound int
		engagement  float64
	}
	byGroupPage := make(map[string]*agg, len(groupsIndex))
	getAgg := func(groupPageID string) *agg {
		if groupPageID == "" {
			return nil
		}
		if a := byGroupPage[groupPageID]; a != nil {
			return a
		}
		a := &agg{}
		byGroupPage[groupPageID] = a
		return a
	}

	// Current period responses.
	respPages, err := nc.QueryAllPages(ctx, dsLogs, map[string]any{
		"filter": map[string]any{
			"and": []any{
				map[string]any{
					"property": "Timestamp",
					"date": map[string]any{
						"on_or_after": curStart.Format(time.RFC3339),
					},
				},
				map[string]any{
					"property": "Direction",
					"select": map[string]any{
						"equals": "Response",
					},
				},
			},
		},
		"page_size": 100,
	})
	if err != nil {
		return fmt.Errorf("query response logs: %w", err)
	}
	for _, p := range respPages {
		groupIDs := notion.RelationIDs(p.Properties, "Group")
		if len(groupIDs) == 0 {
			continue
		}
		a := getAgg(groupIDs[0])
		if a == nil {
			continue
		}
		a.responses++
		if v, ok := notion.NumberValue(p.Properties, "Engagement"); ok {
			a.engagement += v
		}
	}

	// Previous period outbound count (Prev Period Messages).
	outPages, err := nc.QueryAllPages(ctx, dsLogs, map[string]any{
		"filter": map[string]any{
			"and": []any{
				map[string]any{
					"property": "Timestamp",
					"date": map[string]any{
						"on_or_after": prevStart.Format(time.RFC3339),
					},
				},
				map[string]any{
					"property": "Timestamp",
					"date": map[string]any{
						"before": curStart.Format(time.RFC3339),
					},
				},
				map[string]any{
					"property": "Direction",
					"select": map[string]any{
						"equals": "Outbound",
					},
				},
			},
		},
		"page_size": 100,
	})
	if err != nil {
		return fmt.Errorf("query outbound logs: %w", err)
	}
	for _, p := range outPages {
		groupIDs := notion.RelationIDs(p.Properties, "Group")
		if len(groupIDs) == 0 {
			continue
		}
		a := getAgg(groupIDs[0])
		if a == nil {
			continue
		}
		a.prevOutbound++
	}

	// Upsert metric rows for every group we know about.
	for groupPageID, groupID := range groupIDByPageID {
		name := fmt.Sprintf("Group: %s", groupID)
		a := byGroupPage[groupPageID]
		responses := 0
		prevOutbound := 0
		eng := 0.0
		if a != nil {
			responses = a.responses
			prevOutbound = a.prevOutbound
			eng = a.engagement
		}
		props := map[string]any{
			"Name":                notion.Title(name),
			"Group":               notion.Relation(groupPageID),
			"Responses Received":  notion.Number(float64(responses)),
			"Prev Period Messages": notion.Number(float64(prevOutbound)),
			"Total Engagement":    notion.Number(eng),
			"Total Leads":         notion.Number(0),
		}

		if ref, ok := metricsIndex[name]; ok && strings.TrimSpace(ref.ID) != "" {
			if !apply {
				log.Printf("[dry-run] update metrics Name=%q page_id=%s", name, ref.ID)
				continue
			}
			if err := nc.UpdatePageProperties(ctx, ref.ID, props); err != nil {
				return fmt.Errorf("update metrics %q: %w", name, err)
			}
			log.Printf("[update] metrics Name=%q page_id=%s", name, ref.ID)
			continue
		}

		if !apply {
			log.Printf("[dry-run] create metrics Name=%q", name)
			continue
		}
		created, err := nc.CreatePageInDataSource(ctx, dsMetrics, props)
		if err != nil {
			return fmt.Errorf("create metrics %q: %w", name, err)
		}
		log.Printf("[create] metrics Name=%q page_id=%s", name, created.ID)
	}
	return nil
}

func getenvAny(keys ...string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return ""
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func isValidHTTPURL(raw string) bool {
	s := strings.TrimSpace(raw)
	if s == "" {
		return false
	}
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != ""
}
