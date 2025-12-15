package main

import (
	"context"
	"flag"
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
		apply = flag.Bool("apply", false, "write changes to Notion (default: dry-run)")
	)
	flag.Parse()

	// GroupMe token: accept both names (your repo already had GROUP_ME_MASTER_TOKEN).
	gmToken := getenvAny("GROUPME_ACCESS_TOKEN", "GROUP_ME_MASTER_TOKEN")
	if gmToken == "" {
		log.Fatal("missing GroupMe token (set GROUPME_ACCESS_TOKEN or GROUP_ME_MASTER_TOKEN)")
	}

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

	gmc, err := groupme.NewUserClient(groupme.UserConfig{AccessToken: gmToken})
	if err != nil {
		log.Fatal(err)
	}
	nc, err := notion.New(notion.Config{APIKey: notionKey, HTTPTimeout: 60 * time.Second})
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
	log.Printf("mode: %s", ternary(*apply, "APPLY", "DRY-RUN"))

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

	log.Printf("done")
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
