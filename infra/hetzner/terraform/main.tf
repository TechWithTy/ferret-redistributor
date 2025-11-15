terraform {
  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "~> 1.44"
    }
  }
}

provider "hcloud" {
  token = var.hcloud_token
}

resource "hcloud_server" "app" {
  count       = var.app_count
  name        = "social-scale-app-${count.index}"
  image       = "ubuntu-24.04"
  server_type = var.server_type
  ssh_keys    = var.ssh_keys
  location    = var.location
  user_data   = file("${path.module}/cloud-init.yaml")
}

resource "hcloud_volume" "postgres" {
  name     = "social-scale-postgres"
  size     = 50
  format   = "ext4"
  location = var.location
}

resource "hcloud_volume_attachment" "postgres" {
  server_id = hcloud_server.app[0].id
  volume_id = hcloud_volume.postgres.id
  automount = true
}

output "app_ips" {
  value = [for s in hcloud_server.app : s.ipv4_address]
}


