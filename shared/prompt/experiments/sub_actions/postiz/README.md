# Postiz MCP Prompts

This directory contains prompts and sub-prompts for integrating Postiz MCP (Model Context Protocol) server with social media management workflows.

## Overview

The Postiz MCP prompts enable automated social media content creation, scheduling, campaign management, and analytics integration directly from AI development environments like Cursor and Claude.

## Structure

### Master Prompt

- **`master.poml`** - Main orchestrator that coordinates all Postiz MCP operations through sub-prompts.

### Sub-Prompts

1. **`postiz_connect.poml`** - Establishes and validates connections to the Postiz MCP server.
2. **`postiz_create_content.poml`** - Generates platform-optimized content from experiment data.
3. **`postiz_schedule.poml`** - Schedules posts across multiple platforms with optimal timing.
4. **`postiz_campaigns.poml`** - Creates and manages coordinated multi-post campaigns.
5. **`postiz_analytics.poml`** - Retrieves and analyzes post performance metrics.
6. **`postiz_integration.poml`** - Integrates Postiz operations with Notion databases and other SaaS tools.

## Usage

### Basic Workflow

1. **Initialize Connection**: Use `postiz_connect.poml` to verify Postiz MCP server is accessible.
2. **Create Content**: Use `postiz_create_content.poml` to generate platform-optimized content.
3. **Schedule Posts**: Use `postiz_schedule.poml` to schedule posts across platforms.
4. **Manage Campaigns**: Use `postiz_campaigns.poml` for coordinated campaign management.
5. **Analyze Performance**: Use `postiz_analytics.poml` to track and analyze post performance.
6. **Integrate Workflows**: Use `postiz_integration.poml` to sync with Notion and other tools.

### Integration with Experiments

The Postiz prompts integrate seamlessly with the experiment workflow:

- Link posts to Experiments database entries
- Connect to Creative Assets and Scripts/Variants
- Update Copy Calendar with scheduled posts
- Track performance in KPI Progress
- Document learnings in Iterations/Actions

## Prerequisites

- Postiz MCP server configured (see `shared/mcp/_docs/postiz-mcp.md`)
- Postiz API key in `shared/mcp/config.private.json`
- Connected social media accounts in Postiz
- Notion databases set up (Experiments, Creative Assets, Copy Calendar, etc.)

## Configuration

Ensure Postiz MCP server is configured in:
- `shared/mcp/config.public.json` (for local instances)
- `shared/mcp/config.private.json` (for remote instances with API keys)

## Examples

### Schedule a Single Post

```
Execute postiz_schedule.poml:
- Platform: Twitter
- Content: "Just completed refactoring the user authentication module. Excited to share the improvements!"
- Schedule: Today at 2:00 PM EST
- Link to Experiment: EXP-001
```

### Create Multi-Platform Campaign

```
Execute postiz_campaigns.poml:
- Campaign: Product Launch Announcement
- Platforms: Twitter, LinkedIn, Facebook
- Duration: 1 week
- Frequency: 2 posts per day
- Link to Experiment: EXP-002
```

### Analyze Campaign Performance

```
Execute postiz_analytics.poml:
- Campaign ID: CAMP-001
- Metrics: Engagement, Reach, Clicks
- Update KPI Progress: Yes
- Generate Report: Yes
```

## Related Documentation

- [Postiz MCP Documentation](../../../../mcp/_docs/postiz-mcp.md)
- [MCP Configuration Guide](../../../../mcp/README.md)
- [Experiments Master Prompt](../master.poml)

## Notes

- All prompts follow the POML (Prompt Orchestration Markup Language) format.
- Prompts are designed to work with the existing experiment workflow structure.
- All Postiz operations maintain full traceability with Notion databases.
- API keys should never be exposed in prompts or logs.

