# 🐶 Husky Git Hooks
Husky Pre-Commit Hook

This repository uses Husky to run language-specific checks before each commit.

What runs
- Go: `go fmt ./...` and, if available, `golangci-lint run`
- Python: `ruff check`, `black --check`, `isort --check-only`, and `mypy` if installed

Auto-detection
- Go checks run only if Go files or the `go/` module are present
- Python checks run only if Python files or common config files exist, or a `python/` folder exists

Windows notes
- Use Git Bash for hooks. Ensure the hook file uses LF endings
- You can enforce LF with a `.gitattributes` entry (see repo root)

Setup
1) Make the hook executable (Git Bash): `chmod +x .husky/pre-commit`
2) Point Git to Husky hooks: `git config --local core.hooksPath .husky`

Tools
- Go: `go` and optionally `golangci-lint`
- Python: `python` and optionally `ruff`, `black`, `isort`, `mypy`

Notes
- If a tool is not installed, the hook prints a skip message and continues
- Any failing tool will block the commit with a clear error message
This directory contains Git hooks managed by [Husky](https://github.com/automation-co/husky).

## 🚀 Pre-commit Hook

The pre-commit hook runs the following checks before allowing a commit:

1. Checks if a `go` directory exists
2. If found, it will:
   - Run `go fmt ./...` to format all Go code
   - Run `golangci-lint run` for static code analysis

## 🛠️ Installation

1. Install Husky (if not already installed):
   ```bash
   go install github.com/automation-co/husky@latest
   ```

2. Make the pre-commit hook executable (run in Git Bash on Windows):
   ```bash
   chmod +x .husky/pre-commit
   ```

3. If you're on Windows, you might need to configure Git to use the hook:
   ```bash
   git config --local core.hooksPath .husky
   ```

## 🔧 Troubleshooting

- If the hook doesn't run, ensure it's executable:
  ```bash
  ls -la .husky/  # Should show -rwxr-xr-x for pre-commit
  ```

- To bypass the pre-commit hook (use with caution):
  ```bash
  git commit --no-verify -m "Your commit message"
  ```

## ➕ Adding More Hooks

To add more hooks, create new files in this directory with the hook name (e.g., `pre-push`).

## 📝 License

MIT
