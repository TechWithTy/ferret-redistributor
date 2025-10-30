# Husky Hooks

This directory contains Git hooks managed by [Husky](https://github.com/automation-co/husky).

## Pre-commit Hook

The pre-commit hook runs the following checks before allowing a commit:

1. `go fmt ./...` - Formats all Go code
2. `golangci-lint run` - Runs the golangci-lint linter

## Installation

1. Install Husky:
   ```bash
   go install github.com/automation-co/husky@latest
   ```

2. Make the pre-commit hook executable (Unix-like systems):
   ```bash
   chmod +x .husky/pre-commit
   ```

## Adding More Hooks

To add more hooks, create new files in this directory with the hook name (e.g., `pre-push`).
