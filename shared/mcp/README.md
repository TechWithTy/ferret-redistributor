# MCP Configuration

This directory contains cross-language MCP (Model Context Protocol) server configurations that can be used across different backend implementations (Go, Python, TypeScript, Rust, etc.).

## Configuration Files

### `config.public.json`
Public/external MCP server configurations. These are for third-party or publicly available MCP servers that can be shared openly.

**Example:**
```json
{
  "mcpServers": {
    "facebook": {
      "command": "uv",
      "args": [
        "--directory",
        "/path/to/facebook-mcp-server",
        "run",
        "facebook-mcp-server"
      ]
    }
  }
}
```

### `config.private.json`
Private/internal MCP server configurations. These are for internal or proprietary MCP servers that should not be exposed publicly. 

**Security Note:** 
- `config.private.json` is excluded from version control (see `.gitignore`)
- Use `config.private.json.example` as a template
- Copy the example file to create your private config: `cp config.private.json.example config.private.json`
- Use environment variables for sensitive paths, API keys, and credentials

## Configuration Schema

MCP servers can be configured in two ways:

### 1. Command-based MCP Servers

For local or executable MCP servers:

```json
{
  "mcpServers": {
    "server-name": {
      "command": "executable-or-wrapper",
      "args": [
        "argument1",
        "argument2"
      ],
      "env": {
        "ENV_VAR": "value"
      }
    }
  }
}
```

**Fields:**
- **`command`**: The executable command or wrapper (e.g., `uv`, `node`, `python`, `go run`)
- **`args`**: Array of arguments passed to the command
- **`env`** (optional): Environment variables specific to this MCP server

### 2. URL-based MCP Servers

For remote MCP servers accessed via HTTP/HTTPS:

```json
{
  "mcpServers": {
    "server-name": {
      "url": "https://mcp.example.com/server-endpoint"
    }
  }
}
```

**Fields:**
- **`url`**: The HTTP/HTTPS endpoint URL for the remote MCP server
- **`url`** (with token): For direct token authentication, include token as query parameter: `"url": "https://mcp.example.com/server?token=YOUR_TOKEN"` (use in `config.private.json` only)

## Usage Across Languages

### Go
```go
// Load and parse config.public.json or config.private.json
config, err := loadMCPConfig("shared/mcp/config.public.json")
```

### Python
```python
# Load and parse config.public.json or config.private.json
import json
with open('shared/mcp/config.public.json') as f:
    config = json.load(f)
```

### TypeScript
```typescript
// Load and parse config.public.json or config.private.json
import config from './shared/mcp/config.public.json';
```

## Path Resolution

When using these configurations, you may need to resolve relative paths based on your runtime environment. Consider:

1. **Absolute paths**: Use full system paths (e.g., `/path/to/facebook-mcp-server`)
2. **Relative paths**: Resolve relative to the project root or config file location
3. **Environment variables**: Use `$HOME`, `$PROJECT_ROOT`, etc. for portability

## Adding New MCP Servers

1. Determine if the server is public or private
2. Add the configuration to the appropriate file (`config.public.json` or `config.private.json`)
3. Use a descriptive name as the key
4. Document any special requirements or setup steps

## Examples

### Python-based MCP Server
```json
{
  "python-mcp-server": {
    "command": "python",
    "args": [
      "-m",
      "mcp_server_module"
    ],
    "env": {
      "PYTHONPATH": "/path/to/mcp-servers"
    }
  }
}
```

### Node.js-based MCP Server
```json
{
  "node-mcp-server": {
    "command": "node",
    "args": [
      "/path/to/mcp-server/index.js"
    ],
    "env": {
      "NODE_ENV": "production"
    }
  }
}
```

### Go-based MCP Server
```json
{
  "go-mcp-server": {
    "command": "/path/to/compiled-binary",
    "args": [
      "--config",
      "/path/to/config.yaml"
    ]
  }
}
```

### Remote/URL-based MCP Server (Public)
```json
{
  "meta-ads-remote": {
    "url": "https://mcp.pipeboard.co/meta-ads-mcp"
  }
}
```

### Remote/URL-based MCP Server (Private with Token)
For servers requiring authentication, use `config.private.json`:

```json
{
  "meta-ads-remote": {
    "url": "https://mcp.pipeboard.co/meta-ads-mcp?token=YOUR_PIPEBOARD_TOKEN"
  }
}
```

**Note:** Always keep tokens and API keys in `config.private.json`, never in `config.public.json`.

