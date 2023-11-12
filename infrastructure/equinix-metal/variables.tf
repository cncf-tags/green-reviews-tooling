variable "billing_cycle" {
  description = "Billing cycle for the Equinix Metal device"
  type        = string
  default     = "hourly"
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

variable "worker_nodes" {
  description = "List of worker node names"
  type        = list(string)
  default     = ["worker1"]
}
