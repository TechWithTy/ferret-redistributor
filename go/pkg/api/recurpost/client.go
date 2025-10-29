package recurpost

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL   string
	token     string
	http      *http.Client
	userAgent string

	// Services
	Auth      *AuthService
	Accounts  *AccountsService
	Posts     *PostsService
	Schedules *SchedulesService
	Media     *MediaService
    Analytics *AnalyticsService
    UserLogin *UserLoginService
    SocialConnect *SocialConnectService
    SocialAccounts *SocialAccountsService
    Libraries *LibrariesService
}

type Option func(*Client)

func WithBaseURL(u string) Option          { return func(c *Client) { c.baseURL = strings.TrimRight(u, "/") } }
func WithToken(t string) Option            { return func(c *Client) { c.token = t } }
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }
func WithUserAgent(ua string) Option       { return func(c *Client) { c.userAgent = ua } }

func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:   "https://api.recurpost.example",
		http:      &http.Client{Timeout: 30 * time.Second},
		userAgent: "Ferret-RecurPost/0.1",
	}
	for _, o := range opts {
		o(c)
	}
	// wire services
	c.Auth = &AuthService{c}
	c.Accounts = &AccountsService{c}
	c.Posts = &PostsService{c}
	c.Schedules = &SchedulesService{c}
	c.Media = &MediaService{c}
    c.Analytics = &AnalyticsService{c}
    c.UserLogin = &UserLoginService{c}
    c.SocialConnect = &SocialConnectService{c}
    c.SocialAccounts = &SocialAccountsService{c}
    c.Libraries = &LibrariesService{c}
    return c
}

func (c *Client) newRequest(ctx context.Context, method, p string, body io.Reader) (*http.Request, error) {
	if _, err := url.Parse(c.baseURL); err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}
	u := c.baseURL + p
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		var e struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		_ = json.Unmarshal(b, &e)
		return &APIError{StatusCode: resp.StatusCode, Code: e.Code, Message: e.Message}
	}
	if out != nil {
		if err := json.Unmarshal(b, out); err != nil {
			return err
		}
	}
	return nil
}
