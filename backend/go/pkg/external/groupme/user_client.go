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
	u := c.baseURL + "/groups"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, withToken(u, c.token), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("groupme: list groups: %s", res.Status)
	}
	var out apiEnvelope[[]Group]
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	// Normalize members_count if missing
	for i := range out.Response {
		if out.Response[i].MembersCount == 0 && len(out.Response[i].Members) > 0 {
			out.Response[i].MembersCount = len(out.Response[i].Members)
		}
	}
	return out.Response, nil
}

func (c *UserClient) ListBots(ctx context.Context) ([]Bot, error) {
	u := c.baseURL + "/bots"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, withToken(u, c.token), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("groupme: list bots: %s", res.Status)
	}
	var out apiEnvelope[[]Bot]
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Response, nil
}

func withToken(rawURL, token string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := u.Query()
	if q.Get("token") == "" {
		q.Set("token", token)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
