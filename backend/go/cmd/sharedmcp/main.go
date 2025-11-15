package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitesinbyte/ferret/pkg/mcp/sharedcontent"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var (
		flagSharedRoot = flag.String("shared-root", "", "path to the shared/ directory (defaults to ../shared)")
		flagIgnore     = flag.String("ignore", "", "comma-separated list of extra substrings to ignore")
	)
	flag.Parse()

	sharedRoot := firstNonEmpty(
		*flagSharedRoot,
		os.Getenv("SHARED_ASSETS_ROOT"),
		filepath.Clean(filepath.Join("..", "shared")),
	)
	absRoot, err := filepath.Abs(sharedRoot)
	if err != nil {
		log.Fatalf("resolve shared root: %v", err)
	}

	ignoreTokens := sharedcontent.DefaultIgnoreTokens
	if env := os.Getenv("MCP_IGNORE_TOKENS"); env != "" {
		ignoreTokens = append(ignoreTokens, splitTokens(env)...)
	}
	if *flagIgnore != "" {
		ignoreTokens = append(ignoreTokens, splitTokens(*flagIgnore)...)
	}

	srv := server.NewMCPServer(
		"Shared Assets",
		"0.1.0",
		server.WithResourceCapabilities(false, true),
		server.WithResourceRecovery(),
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	if err := sharedcontent.RegisterSharedContent(
		srv,
		absRoot,
		sharedcontent.WithIgnoreTokens(ignoreTokens),
	); err != nil {
		log.Fatalf("register shared content: %v", err)
	}

	log.Printf("shared MCP server exposing %s (ignore tokens: %v)", absRoot, ignoreTokens)

	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("shared MCP server error: %v", err)
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func splitTokens(raw string) []string {
	chunks := strings.Split(raw, ",")
	var tokens []string
	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}
		tokens = append(tokens, chunk)
	}
	return tokens
}

