terraform {
  required_providers {
    equinix = {
      source  = "equinix/equinix"
      version = "1.13.0"
    }
    null = {
      source = "hashicorp/null"
      version = "3.2.2"
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

  connection {
    user = "root"
    private_key = file(var.ssh_private_key_path)
    host = self.access_public_ipv4
  }

  provisioner "remote-exec" {
    inline = [
      "curl -sfL https://get.k3s.io | INSTALL_K3S_CHANNEL=${var.k3s_version} K3S_TOKEN=${var.k3s_agent_token} sh -s - server --node-taint CriticalAddonsOnly=true:NoExecute --flannel-backend=none --disable-network-policy",
      "systemctl is-active --quiet k3s.service",
    ]
  }
}

// NOTE: to extract KUBECONFIG we need to copy kubeconfig in controlplane:/etc/rancher/k3s/k3s.yaml

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

resource "null_resource" "install_cilium_cni" {
  depends_on          = [equinix_metal_device.control_plane]

  connection {
    user = "root"
    private_key = file(var.ssh_private_key_path)
    host = equinix_metal_device.control_plane.access_public_ipv4
  }

  provisioner "remote-exec" {
    inline = [
      "echo '@@@@@@ Installing Cilium @@@@@@'",
      "CILIUM_CLI_VERSION=$(curl -s https://raw.githubusercontent.com/cilium/cilium-cli/main/stable.txt)",
      "CLI_ARCH=amd64",
      "curl -L --fail --remote-name-all https://github.com/cilium/cilium-cli/releases/download/$${CILIUM_CLI_VERSION}/cilium-linux-$${CLI_ARCH}.tar.gz{,.sha256sum}",
      "sha256sum --check cilium-linux-$${CLI_ARCH}.tar.gz.sha256sum",
      "sudo tar xzvfC cilium-linux-$${CLI_ARCH}.tar.gz /usr/local/bin",
      "rm cilium-linux-$${CLI_ARCH}.tar.gz{,.sha256sum}",
      "echo '@@@@@@ Installed Cilium @@@@@@'",
      "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml",
      "echo '@@@@@@ Adding Cilium CNI to cluster @@@@@@'",
      "cilium install --version ${var.cilium_version}",
      "cilium status --wait"
    ]
  }
}

resource "null_resource" "bootstrap_flux" {
  depends_on          = [null_resource.install_cilium_cni]

  connection {
    user = "root"
    private_key = file(var.ssh_private_key_path)
    host = equinix_metal_device.control_plane.access_public_ipv4
  }

  provisioner "remote-exec" {
    inline = [
      "curl -s https://fluxcd.io/install.sh | sudo FLUX_VERSION=${var.flux_version} bash",
      "export GITHUB_TOKEN=${var.flux_github_token}",
      "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml",
      "flux bootstrap github --owner=${var.flux_github_user} --repository=green-reviews-tooling --path=clusters"
    ]
  }
}
