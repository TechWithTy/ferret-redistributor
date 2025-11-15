# CI/CD & Deployment Secrets

## Required GitHub Secrets / Environments

| Secret | Description | Used by |
| ------ | ----------- | ------- |
| `GHCR_PAT` *(optional)* | If pushing to GHCR from non-default context. Otherwise the repo `GITHUB_TOKEN` suffices. | `.github/workflows/build.yml` |
| `HETZNER_SSH_KEY` | Private key with access to the Hetzner nodes. Stored as an Actions secret scoped to the `hetzner-prod` environment. | `.github/workflows/deploy-hetzner.yml` |
| `HETZNER_HOST` | Public IP / hostname of the primary Hetzner host running the compose stack. | `.github/workflows/deploy-hetzner.yml` |
| `HCLOUD_TOKEN` *(optional)* | For running Terraform in CI if you automate infra provisioning. | Future infra workflows |

## Deployment Flow

1. CI runs `lint.yml` + `build.yml` automatically on pushes / PRs.
2. When targeting `main`, the `deploy-hetzner` job executes `shared/tools/deploy_hetzner.sh`, which:
   - SSHes into the Hetzner box
   - Pulls the latest commit
   - Runs `docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d`
3. Ensure the compose stack can reach the secrets mounted on the host (e.g., `/etc/social-scale/env` or Docker secrets).

Before enabling the deployment job, add the secrets above in **Settings → Secrets and variables → Actions**, and protect the `hetzner-prod` environment with required reviewers if needed.


