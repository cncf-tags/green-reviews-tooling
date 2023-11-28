terraform {
  required_providers {
    equinix = {
      source  = "equinix/equinix"
      version = "1.13.0"
    }
  }

  backend "s3" {
    bucket  = "green-reviews-state-bucket"
    key     = "opentofu/terraform.tfstate"
    region  = "eu-central-1"
    encrypt = true
  }
}

provider "equinix" {
  auth_token = var.equinix_auth_token
}

resource "equinix_metal_project_ssh_key" "ssh_key" {
  name       = var.cluster_name
  project_id = var.project_id
  public_key = var.ssh_public_key
}

resource "equinix_metal_device" "control_plane" {
  hostname            = "${var.cluster_name}-control-plane"
  plan                = var.device_plan
  metro               = var.device_metro
  operating_system    = var.device_os
  billing_cycle       = var.billing_cycle
  project_id          = var.project_id
  depends_on          = [equinix_metal_project_ssh_key.ssh_key]
  project_ssh_key_ids = [equinix_metal_project_ssh_key.ssh_key.id]
  user_data           = <<EOF
  #!/bin/bash
  curl -sfL https://get.k3s.io | INSTALL_K3S_CHANNEL="${var.k3s_version}" K3S_TOKEN="${var.k3s_agent_token}" sh -s - server \
	--node-taint CriticalAddonsOnly=true:NoExecute \
	--flannel-backend=none \
	--disable-network-policy
  EOF


  behavior {
    allow_changes = [
      "user_data"
    ]
  }
}

// NOTE: need to check if we need to execute a command in the controlplane to verify that k3s.service is up and running
// before we try to make workerplanes to join

// to extract KUBECONFIG we need to copy file in controlplane:/etc/rancher/k3s/k3s.yaml

resource "equinix_metal_device" "worker" {
  for_each            = toset(var.worker_nodes)
  hostname            = "${var.cluster_name}-${each.value}"
  plan                = var.device_plan
  metro               = var.device_metro
  operating_system    = var.device_os
  billing_cycle       = var.billing_cycle
  project_id          = var.project_id
  project_ssh_key_ids = [equinix_metal_project_ssh_key.ssh_key.id]
  depends_on          = [equinix_metal_device.control_plane]
  user_data           = <<EOF
  #!/bin/bash
  curl -sfL https://get.k3s.io | INSTALL_K3S_CHANNEL="${var.k3s_version}" sh -s - agent --token "${var.k3s_agent_token}" --server "https://${equinix_metal_device.control_plane.access_private_ipv4}:6443"
  EOF

  behavior {
    allow_changes = [
      "user_data"
    ]
  }
}
