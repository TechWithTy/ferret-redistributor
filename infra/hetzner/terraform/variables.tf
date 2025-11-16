variable "hcloud_token" {
  type        = string
  description = "Hetzner API token"
}

variable "ssh_keys" {
  type        = list(string)
  description = "Names of SSH keys uploaded to Hetzner"
}

variable "location" {
  type        = string
  default     = "nbg1"
  description = "Hetzner location"
}

variable "server_type" {
  type        = string
  default     = "cpx31"
}

variable "app_count" {
  type    = number
  default = 1
}




