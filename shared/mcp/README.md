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

## Postiz MCP Server

Postiz is an open-source social media scheduling tool that provides an MCP server for automating social media content creation and scheduling directly from AI development environments like Cursor and Claude.

**📖 For detailed documentation, see:** [`_docs/postiz-mcp.md`](_docs/postiz-mcp.md)

### Overview

The Postiz MCP server enables you to:
- Schedule social media posts across multiple platforms (Twitter, LinkedIn, Facebook, Instagram, etc.)
- Create and manage content directly from your coding environment
- Automate content workflows by integrating with other SaaS applications
- Manage multiple social media accounts from a single interface

### Setup Options

#### Option 1: Self-Hosted Postiz MCP Server

If you're running your own Postiz instance, you can set up the MCP server locally:

1. **Using Docker Compose:**
   ```bash
   # Clone or set up Postiz MCP server
   # Configure .env file with your Postiz API credentials
   docker-compose up -d
   ```

2. **Using Docker Directly:**
   ```bash
   docker build -t oculair/postiz-mcp:latest .
   docker run -p 3084:3084 --env-file .env --rm -it oculair/postiz-mcp:latest
   ```

3. **Environment Variables:**
   ```bash
   POSTIZ_API_URL=https://your-postiz-instance.com/api
   POSTIZ_API_KEY=your_postiz_api_key_here
   PORT=3084
   NODE_ENV=production
   ```

#### Option 2: Remote Postiz MCP Server

If using a hosted Postiz service, configure the remote endpoint directly.

### Configuration

#### For Self-Hosted Instances

Add to `config.private.json`:

```json
{
  "mcpServers": {
    "postiz": {
      "url": "http://localhost:3084/sse"
    }
  }
}
```

#### For Remote Instances with API Key

Add to `config.private.json`:

```json
{
  "mcpServers": {
    "postiz": {
      "url": "https://your-postiz-instance.com/mcp/sse?api_key=YOUR_API_KEY"
    }
  }
}
```

**Security Note:** Always use `config.private.json` for Postiz configuration as it requires API keys. Never commit API keys to version control.

### Usage Examples

Once configured, you can use Postiz MCP from your AI assistant:

**Example 1: Schedule a Single Post**
```
Schedule a post on Twitter about completing the refactoring of the user authentication module.
```

**Example 2: Multi-Platform Posting**
```
Post about the latest bug fix on Twitter, LinkedIn, and Facebook.
```

**Example 3: Detailed Post with Hashtags**
```
Schedule a post on LinkedIn about our new feature launch. Include relevant hashtags and mention it's available now.
```

### Features

- **Multi-Platform Support**: Schedule posts across Twitter, LinkedIn, Facebook, Instagram, and more
- **Content Templates**: Create reusable templates for consistent branding
- **Scheduling**: Schedule posts in advance with optimal timing
- **Analytics**: Track post performance and engagement metrics
- **Team Collaboration**: Manage team permissions and workflows
- **SaaS Integration**: Connect with Notion, CRM systems, and other tools for automated content workflows

### Integration with Development Workflows

Postiz MCP integrates seamlessly with:
- **Cursor**: Schedule posts directly from your coding environment
- **Claude Desktop**: Manage social media from your AI assistant
- **Notion**: Automate content from your project documentation
- **Other SaaS Tools**: Connect to CRM, analytics, and project management platforms

### Resources

- **Postiz Website**: https://postiz.com
- **Postiz Documentation**: https://postiz.com/docs
- **MCP Protocol**: Learn more about Model Context Protocol

### Troubleshooting

**Connection Issues:**
- Verify the Postiz MCP server is running (check port 3084 for local instances)
- Ensure your API key is correct and has proper permissions
- Check firewall settings if using a remote instance

**Authentication Errors:**
- Verify your `POSTIZ_API_KEY` is set correctly
- Ensure the API key has the necessary permissions for MCP access
- Check that the Postiz instance URL is correct

**Post Scheduling Failures:**
- Verify your social media accounts are connected in Postiz
- Check that the content meets platform-specific requirements
- Ensure scheduling permissions are enabled for your account

