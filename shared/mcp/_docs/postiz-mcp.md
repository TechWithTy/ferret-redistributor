# Postiz MCP Server Documentation

## Overview

Postiz is an open-source social media scheduling tool that provides an MCP (Model Context Protocol) server for automating social media content creation and scheduling directly from AI development environments like Cursor and Claude.

The Postiz MCP server enables seamless integration between your development workflow and social media management, allowing you to schedule posts, manage content, and automate social media workflows without leaving your coding environment.

## Key Features

- **Multi-Platform Support**: Schedule posts across Twitter, LinkedIn, Facebook, Instagram, WhatsApp, and more
- **Content Templates**: Create reusable templates for consistent branding across platforms
- **Advanced Scheduling**: Schedule posts in advance with optimal timing based on audience analytics
- **Performance Analytics**: Track post performance and engagement metrics
- **Team Collaboration**: Manage team permissions and collaborative content workflows
- **SaaS Integration**: Connect with Notion, CRM systems, and other tools for automated content workflows
- **AI-Powered Content**: Leverage AI for content suggestions and optimization

## Architecture

Postiz MCP uses the Model Context Protocol to provide a standardized interface for social media management. The server communicates via Server-Sent Events (SSE) over HTTP/HTTPS, making it compatible with various MCP clients.

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Cursor    │────────▶│  Postiz MCP  │────────▶│   Postiz    │
│   Claude    │  MCP    │    Server    │  API    │   Platform  │
│             │         │              │         │             │
└─────────────┘         └──────────────┘         └─────────────┘
```

## Setup Instructions

### Prerequisites

- Postiz account and API key
- Docker (for self-hosted setup) or access to hosted Postiz instance
- MCP-compatible client (Cursor, Claude Desktop, etc.)

### Option 1: Self-Hosted Postiz MCP Server

#### Using Docker Compose

1. **Clone or download the Postiz MCP server repository**

2. **Create environment file:**
   ```bash
   cp .env.example .env
   ```

3. **Configure environment variables:**
   ```bash
   POSTIZ_API_URL=https://your-postiz-instance.com/api
   POSTIZ_API_KEY=your_postiz_api_key_here
   PORT=3084
   NODE_ENV=production
   ```

4. **Start the server:**
   ```bash
   docker-compose up -d
   ```

#### Using Docker Directly

1. **Build the Docker image:**
   ```bash
   docker build -t oculair/postiz-mcp:latest .
   ```

2. **Run the container:**
   ```bash
   docker run -p 3084:3084 --env-file .env --rm -it oculair/postiz-mcp:latest
   ```

#### Verify Installation

Check that the server is running:
```bash
curl http://localhost:3084/health
```

### Option 2: Remote Postiz MCP Server

If using a hosted Postiz service, configure the remote endpoint directly in your MCP configuration.

## Configuration

### For Cursor

1. Navigate to `Settings` → `Cursor Settings` → `MCP` → `Add new global MCP server`
2. Add the following configuration:

**Local Instance:**
```json
{
  "mcpServers": {
    "postiz": {
      "url": "http://localhost:3084/sse",
      "disabled": false,
      "alwaysAllow": []
    }
  }
}
```

**Remote Instance:**
```json
{
  "mcpServers": {
    "postiz": {
      "url": "https://your-postiz-instance.com/mcp/sse?api_key=YOUR_API_KEY",
      "disabled": false,
      "alwaysAllow": []
    }
  }
}
```

3. Save and restart Cursor

### For Claude Desktop

1. Go to `Settings` → `Connectors` → `Add Custom Connector`
2. Enter `Postiz` as the connector name
3. Paste the Postiz MCP server URL (e.g., `http://localhost:3084/sse`)
4. Click `Add` and restart Claude Desktop

### Project Configuration

Add to `shared/mcp/config.private.json`:

```json
{
  "mcpServers": {
    "postiz": {
      "url": "http://localhost:3084/sse"
    },
    "postiz-remote": {
      "url": "https://your-postiz-instance.com/mcp/sse?api_key=YOUR_POSTIZ_API_KEY"
    }
  }
}
```

**Security Note:** Always use `config.private.json` for Postiz configuration as it requires API keys. Never commit API keys to version control.

## Usage Examples

### Basic Post Scheduling

**Example 1: Single Platform Post**
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

### Advanced Workflows

**Example 4: Scheduled Campaign**
```
Create a series of 5 posts about our product launch, schedule them over the next week on Twitter and LinkedIn, and include images.
```

**Example 5: Content from Notion**
```
Automatically schedule a post about the project milestone we just completed in Notion.
```

**Example 6: Performance-Based Posting**
```
Post about our latest metrics when our user count reaches 10,000.
```

## Integration Workflows

### Development Workflow Integration

Postiz MCP integrates seamlessly with your development workflow:

1. **Code Completion → Social Update**
   - Complete a feature
   - Ask AI assistant to schedule a post about it
   - Post is automatically scheduled across platforms

2. **Project Milestones → Content**
   - Update project status in Notion
   - Trigger automated social media post
   - Keep audience informed automatically

3. **Bug Fixes → Community Updates**
   - Fix a critical bug
   - Schedule a post explaining the fix
   - Build trust with transparent communication

### SaaS Integration

Connect Postiz to various SaaS applications:

- **Notion**: Automate content from project documentation
- **CRM Systems**: Share customer wins and testimonials
- **Analytics Dashboards**: Post about key performance metrics
- **Project Management Tools**: Announce milestones and releases
- **Support Platforms**: Share positive customer feedback

## API Reference

### Endpoints

#### Health Check
```
GET /health
```
Returns server health status.

#### SSE Endpoint
```
GET /sse
```
Server-Sent Events endpoint for MCP communication.

### Authentication

Postiz MCP server requires authentication via:
- API key in query parameter: `?api_key=YOUR_API_KEY`
- Or environment variable: `POSTIZ_API_KEY`

### Rate Limits

Check Postiz API documentation for current rate limits. The MCP server respects these limits and will queue requests as needed.

## Best Practices

### Content Strategy

1. **Maintain Consistency**: Use content templates to ensure brand voice consistency
2. **Platform Optimization**: Adapt content format for each platform's audience
3. **Timing Matters**: Schedule posts at optimal times based on audience analytics
4. **Engage Authentically**: Balance automation with genuine engagement

### Security

1. **API Key Management**: 
   - Store API keys in `config.private.json` only
   - Never commit API keys to version control
   - Rotate keys regularly

2. **Access Control**:
   - Use team permissions in Postiz to control who can schedule posts
   - Review scheduled posts before publication when possible

3. **Content Review**:
   - Set up approval workflows for sensitive content
   - Monitor scheduled posts regularly

### Performance

1. **Batch Operations**: Group multiple posts when possible
2. **Error Handling**: Monitor failed posts and retry as needed
3. **Analytics**: Regularly review post performance to optimize strategy

## Troubleshooting

### Connection Issues

**Problem:** Cannot connect to Postiz MCP server

**Solutions:**
- Verify the Postiz MCP server is running (check port 3084 for local instances)
- Ensure your API key is correct and has proper permissions
- Check firewall settings if using a remote instance
- Verify the URL format is correct: `http://localhost:3084/sse` or `https://your-instance.com/mcp/sse`

### Authentication Errors

**Problem:** Authentication failed

**Solutions:**
- Verify your `POSTIZ_API_KEY` is set correctly
- Ensure the API key has the necessary permissions for MCP access
- Check that the Postiz instance URL is correct
- Verify the API key hasn't expired or been revoked

### Post Scheduling Failures

**Problem:** Posts fail to schedule

**Solutions:**
- Verify your social media accounts are connected in Postiz
- Check that the content meets platform-specific requirements (character limits, image formats, etc.)
- Ensure scheduling permissions are enabled for your account
- Review Postiz logs for specific error messages

### Performance Issues

**Problem:** Slow response times

**Solutions:**
- Check network connectivity
- Verify Postiz API status
- Review rate limiting settings
- Consider using a local instance for faster response times

## Resources

- **Postiz Website**: https://postiz.com
- **Postiz Documentation**: https://postiz.com/docs
- **Postiz GitHub**: https://github.com/gitroomhq/postiz-app
- **MCP Protocol Documentation**: https://modelcontextprotocol.io
- **Postiz Community**: Check Reddit and Product Hunt for community discussions

## Support

For issues specific to:
- **Postiz Platform**: Contact Postiz support or check their documentation
- **MCP Integration**: Review MCP protocol documentation
- **Project Configuration**: Check `shared/mcp/README.md` for configuration details

## Changelog

### Version 1.0.0
- Initial Postiz MCP server documentation
- Basic setup and configuration instructions
- Usage examples and best practices

## Related Documentation

- [MCP Configuration Guide](../README.md)
- [Postiz Integration Plan](../../../app/app/POSTIZ_INTEGRATION_PLAN.md)

