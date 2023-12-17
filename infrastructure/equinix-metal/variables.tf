variable "billing_cycle" {
  description = "Billing cycle for the Equinix Metal device"
  type        = string
  default     = "hourly"
}

variable "cilium_version" {
  description = "cilium version for the cluster"
  type        = string
  default     = "1.14.4"
}

variable "cluster_name" {
  description = "Name of the cluster"
  type        = string
  default     = "green-reviews"
}

variable "device_metro" {
  description = "Metro location for the Equinix Metal device"
  type        = string
  default     = "pa"
}

variable "device_os" {
  description = "Operating system for the Equinix Metal device"
  type        = string
  default     = "ubuntu_22_04"
}

variable "device_plan" {
  description = "Plan type for the Equinix Metal device"
  type        = string
  default     = "m3.small.x86"
}

variable "flux_github_token" {
  description = "GitHub token for Flux"
  type        = string
  sensitive = true
}

variable "flux_github_user" {
  description = "GitHub user for Flux"
  type        = string
  default     = "cncf-tags"
}

variable "flux_version" {
  description = "Flux cli version"
  type        = string
  default     = "2.1.2"
}

variable "k3s_version" {
  description = "k3s version for the cluster"
  type        = string
  default     = "v1.28.3+k3s2"
}

variable "k3s_agent_token" {
  description = "agent token for joining workernodes with the controlplane"
  type = string
  sensitive = true
}

variable "equinix_auth_token" {
  description = "Authentication token for Equinix Metal"
  type        = string
  sensitive   = true
}

variable "project_id" {
  description = "Project ID for the Equinix Metal resources"
  type        = string
  sensitive   = true
}

variable "ssh_public_key" {
  description = "SSH public key for the Equinix Metal device"
  type        = string
  sensitive   = true
}


variable "ssh_private_key_path" {
  description = "SSH private key path for the Equinix Metal device"
  type        = string
  default = "~/.ssh/id_rsa"
}

variable "worker_nodes" {
  description = "List of worker node names"
  type        = list(string)
  default     = ["worker1"]
}
