# Hetzner Terraform Stack

Usage:

```bash
cd infra/hetzner/terraform
terraform init
terraform plan \
  -var "hcloud_token=..." \
  -var 'ssh_keys=["laptop-key"]'
terraform apply ...
```

Resources:

- `hcloud_server.app` – Ubuntu hosts that run the Docker compose stack.
- `hcloud_volume.postgres` – Persistent volume for database state.
- `cloud-init.yaml` – Installs Docker + systemd unit that runs `docker compose up` on boot.

The GitHub Actions deployment job will SSH into these hosts and run
`shared/tools/deploy_hetzner.sh`, reusing the same compose files as local dev.




