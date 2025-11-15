package sharedcontent

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// DefaultIgnoreTokens defines substrings that, when present in a path, prevent the
// asset from being exposed through the MCP server.
var DefaultIgnoreTokens = []string{
	"_service",
	"_auth",
	"_mcp_deny",
}

// Option configures RegisterSharedContent behavior.
type Option func(*config)

// WithIgnoreTokens overrides the default set of sensitive substrings.
func WithIgnoreTokens(tokens []string) Option {
	return func(cfg *config) {
		if len(tokens) == 0 {
			return
		}
		cfg.ignoreTokens = tokens
	}
}

type config struct {
	ignoreTokens []string
}

// RegisterSharedContent walks the shared assets directory (excluding the generated
// subtree) and registers each file as an MCP resource. POML files are additionally
// exposed as prompts.
func RegisterSharedContent(
	srv *server.MCPServer,
	sharedRoot string,
	opts ...Option,
) error {
	info, err := os.Stat(sharedRoot)
	if err != nil {
		return fmt.Errorf("sharedcontent: stat shared root: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("sharedcontent: %s is not a directory", sharedRoot)
	}

	cfg := config{
		ignoreTokens: DefaultIgnoreTokens,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	assets, err := discoverAssets(sharedRoot, cfg.ignoreTokens)
	if err != nil {
		return err
	}

	for _, asset := range assets {
		if err := registerResource(srv, asset); err != nil {
			return err
		}
		if asset.category == categoryPrompt {
			if err := registerPrompt(srv, asset); err != nil {
				return err
			}
		}
	}

	return nil
}

type assetCategory int

const (
	categoryResource assetCategory = iota
	categoryTool
	categoryPrompt
)

func (c assetCategory) String() string {
	switch c {
	case categoryTool:
		return "tool"
	case categoryPrompt:
		return "prompt"
	default:
		return "resource"
	}
}

type asset struct {
	absPath  string
	relPath  string
	category assetCategory
}

func discoverAssets(root string, ignoreTokens []string) ([]asset, error) {
	var assets []asset
	lowerTokens := normalizeTokens(ignoreTokens)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if path == root {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		if relPath == "generated" || strings.HasPrefix(relPath, "generated/") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldIgnore(relPath, d.IsDir(), lowerTokens) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		category := categorize(relPath)
		assets = append(assets, asset{
			absPath:  path,
			relPath:  relPath,
			category: category,
		})
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("sharedcontent: walk shared directory: %w", err)
	}

	return assets, nil
}

func normalizeTokens(tokens []string) []string {
	var out []string
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		out = append(out, strings.ToLower(token))
	}
	return out
}

func shouldIgnore(relPath string, isDir bool, tokens []string) bool {
	lower := strings.ToLower(relPath)
	for _, token := range tokens {
		if token == "" {
			continue
		}
		if strings.Contains(lower, token) {
			return true
		}
	}
	// dotfiles are ignored by default
	if base := filepath.Base(relPath); strings.HasPrefix(base, ".") && base != ".gitkeep" {
		return true
	}
	if !isDir && strings.HasSuffix(lower, "~") {
		return true
	}
	return false
}

func categorize(relPath string) assetCategory {
	segments := strings.Split(relPath, "/")
	if len(segments) == 0 {
		return categoryResource
	}
	top := segments[0]
	ext := strings.ToLower(filepath.Ext(relPath))

	switch {
	case top == "prompt" || top == "prompts":
		if ext == ".poml" {
			return categoryPrompt
		}
		return categoryResource
	case top == "tools":
		return categoryTool
	case ext == ".poml":
		return categoryPrompt
	default:
		return categoryResource
	}
}

func registerResource(srv *server.MCPServer, asset asset) error {
	uri := fmt.Sprintf("shared://%s", asset.relPath)
	title := fmt.Sprintf("%s (%s)", filepath.Base(asset.relPath), asset.category)
	mime := mimeTypeFor(asset.relPath)

	resource := mcp.NewResource(
		uri,
		title,
		mcp.WithMIMEType(mime),
		mcp.WithResourceDescription(fmt.Sprintf("Shared %s file at %s", asset.category, asset.relPath)),
	)

	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		bytes, err := os.ReadFile(asset.absPath)
		if err != nil {
			return nil, err
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: mime,
				Text:     string(bytes),
			},
		}, nil
	}

	srv.AddResource(resource, handler)
	return nil
}

func registerPrompt(srv *server.MCPServer, asset asset) error {
	name := promptName(asset.relPath)
	prompt := mcp.NewPrompt(
		name,
		mcp.WithPromptDescription(fmt.Sprintf("Shared prompt sourced from %s", asset.relPath)),
		mcp.WithArgument(
			"context",
			mcp.ArgumentDescription("Optional supplemental context appended to the prompt."),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		body, err := os.ReadFile(asset.absPath)
		if err != nil {
			return nil, err
		}

		content := string(body)
		if extra, ok := request.Params.Arguments["context"]; ok && extra != "" {
			content = content + "\n\n" + extra
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Shared prompt %s", name),
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent(content),
				),
			},
		), nil
	}

	srv.AddPrompt(prompt, handler)
	return nil
}

func promptName(relPath string) string {
	name := strings.TrimSuffix(relPath, filepath.Ext(relPath))
	name = strings.ReplaceAll(name, "/", "_")
	return name
}

func mimeTypeFor(relPath string) string {
	switch strings.ToLower(filepath.Ext(relPath)) {
	case ".md":
		return "text/markdown"
	case ".json":
		return "application/json"
	case ".yaml", ".yml":
		return "application/yaml"
	case ".poml":
		return "text/poml"
	case ".py":
		return "text/x-python"
	case ".sh":
		return "text/x-shellscript"
	case ".toon":
		return "text/vnd.socialscale.toon"
	default:
		return "text/plain"
	}
}
