# 🐶 Husky Git Hooks

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
