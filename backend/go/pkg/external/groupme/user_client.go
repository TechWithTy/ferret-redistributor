package groupme

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// UserClient calls GroupMe endpoints that require a user access token (not a bot id).
type UserClient struct {
	baseURL string
	token   string
	hc      *http.Client
}

type UserConfig struct {
	BaseURL     string
	AccessToken string
	HTTPTimeout time.Duration
}

// NewUserClient creates a client for listing groups and bots via the GroupMe API.
func NewUserClient(cfg UserConfig) (*UserClient, error) {
	tok := strings.TrimSpace(cfg.AccessToken)
	if tok == "" {
		return nil, errors.New("groupme: missing access token")
	}
	base := strings.TrimSpace(cfg.BaseURL)
	if base == "" {
		base = defaultBaseURL
	}
	to := cfg.HTTPTimeout
	if to <= 0 {
		to = 15 * time.Second
	}
	return &UserClient{
		baseURL: strings.TrimRight(base, "/"),
		token:   tok,
		hc:      &http.Client{Timeout: to},
	}, nil
}

// Group is a minimal GroupMe group object.
type Group struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CreatorUserID string `json:"creator_user_id"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	MembersCount  int    `json:"members_count"`
	// Some responses only have "members": [...]
	Members []struct {
		UserID string `json:"user_id"`
	} `json:"members"`
}

// Bot is a minimal GroupMe bot object.
type Bot struct {
	BotID       string `json:"bot_id"`
	GroupID     string `json:"group_id"`
	Name        string `json:"name"`
	AvatarURL   string `json:"avatar_url"`
	CallbackURL string `json:"callback_url"`
	CreatedAt   int64  `json:"created_at"`
}

type apiEnvelope[T any] struct {
	Response T `json:"response"`
}

func (c *UserClient) ListGroups(ctx context.Context) ([]Group, error) {
	items, err := listPaged[Group](ctx, c.hc, c.baseURL, c.token, "/groups")
	if err != nil {
		return nil, err
	}
	// Normalize members_count if missing
	for i := range items {
		if items[i].MembersCount == 0 && len(items[i].Members) > 0 {
			items[i].MembersCount = len(items[i].Members)
		}
	}
	return items, nil
}

func (c *UserClient) ListBots(ctx context.Context) ([]Bot, error) {
	return listPaged[Bot](ctx, c.hc, c.baseURL, c.token, "/bots")
}

func listPaged[T any](ctx context.Context, hc *http.Client, baseURL, token, path string) ([]T, error) {
	const perPage = 100
	out := make([]T, 0, 128)
	for page := 1; page < 10_000; page++ { // safety cap
		u := strings.TrimRight(baseURL, "/") + path
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, withTokenAndPaging(u, token, page, perPage), nil)
		if err != nil {
			return nil, err
		}
		res, err := hc.Do(req)
		if err != nil {
			return nil, err
		}
		var env apiEnvelope[[]T]
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			_ = res.Body.Close()
			return nil, fmt.Errorf("groupme: %s: %s", strings.TrimPrefix(path, "/"), res.Status)
		}
		if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
			_ = res.Body.Close()
			return nil, err
		}
		_ = res.Body.Close()

		if len(env.Response) == 0 {
			break
		}
		out = append(out, env.Response...)
	}
	return out, nil
}

func withTokenAndPaging(rawURL, token string, page, perPage int) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := u.Query()
	if q.Get("token") == "" {
		q.Set("token", token)
	}
	if q.Get("page") == "" && page > 0 {
		q.Set("page", fmt.Sprintf("%d", page))
	}
	if q.Get("per_page") == "" && perPage > 0 {
		q.Set("per_page", fmt.Sprintf("%d", perPage))
	}
	u.RawQuery = q.Encode()
	return u.String()
}
