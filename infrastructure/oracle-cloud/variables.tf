variable "api_key_fingerprint" {
  description = "Fingerprint of the API key"
  type        = string
}

variable "api_private_key_path" {
  description = "Path to the private key"
  type        = string
}

variable "availability_domain" {
  description = "The availability domain to create resources in"
  type        = string
  default     = "vncY:US-ASHBURN-AD-1"
}

variable "bastion_subnet_cidr" {
  description = "CIDR block for the subnet"
  type        = string
  default     = "10.0.1.0/24"
}

variable "benchmark_subnet_cidr" {
  description = "CIDR block for the subnet"
  type        = string
  default     = "10.0.2.0/24"
}

variable "bm_shape" {
  description = "Shape of the BM instance"
  type        = string
  default     = "BM.Standard2.52"
}

variable "bm_name" {
  description = "Display name for the BM instance"
  type        = string
  default     = "green-reviews-benchmark"
}

variable "compartment_ocid" {
  description = "OCID of the compartment"
  type        = string
}

variable "k3s_version" {
  description = "k3s version for the cluster"
  type        = string
  default     = "v1.31.6+k3s1"
}

variable "operating_system" {
  description = "Operating system for the BM instance"
  type        = string
  default     = "Canonical Ubuntu"
}

variable "operating_system_version" {
  description = "Operating system version"
  type        = string
  default     = "24.04"
}

variable "region" {
  description = "Region for resources to be created"
  type        = string
  default     = "us-ashburn-1"
}

variable "ssh_public_key" {
  description = "SSH public key for instance access"
  type        = string
}

variable "ssh_private_key_path" {
  description = "SSH private key path for instance access"
  type        = string
  default     = "~/.ssh/id_rsa"
}

variable "ssh_user" {
  description = "SSH user for the instances"
  type        = string
  default     = "ubuntu"
}

variable "tenancy_ocid" {
  description = "OCID of your tenancy"
  type        = string
}

variable "user_ocid" {
  description = "OCID of the user"
  type        = string
}

variable "vcn_cidr" {
  description = "CIDR block for the VCN"
  type        = string
  default     = "10.0.0.0/16"
}

variable "vcn_name" {
  description = "Display name for the VCN"
  type        = string
  default     = "green-reviews-vcn"
}

variable "vm_ocpus" {
  description = "Number of OCPUs for the instance"
  type        = number
  default     = 1
}

variable "vm_memory_in_gbs" {
  description = "Amount of memory in GBs for the instance"
  type        = number
  default     = 2
}

variable "vm_shape" {
  description = "Shape of the VM instance"
  type        = string
  default     = "VM.Standard.E4.Flex"
}

variable "vm_name" {
  description = "Display name for the VM instance"
  type        = string
  default     = "green-reviews-bastion"
}
