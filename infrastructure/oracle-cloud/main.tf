terraform {
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "6.29.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "3.2.2"
    }
  }
}

provider "oci" {
  tenancy_ocid     = var.tenancy_ocid
  user_ocid        = var.user_ocid
  fingerprint      = var.api_key_fingerprint
  private_key_path = var.api_private_key_path
  region           = var.region
}

# Get the latest Ubuntu 24.04 image.
data "oci_core_images" "ubuntu_images" {
  compartment_id           = var.compartment_ocid
  operating_system         = var.operating_system
  operating_system_version = var.operating_system_version
  shape                    = var.bm_shape
  sort_by                  = "TIMECREATED"
  sort_order               = "DESC"
}

resource "oci_core_vcn" "gr_vcn" {
  compartment_id = var.compartment_ocid
  cidr_block     = var.vcn_cidr
  display_name   = var.vcn_name
  dns_label      = "grvcn"
}

resource "oci_core_internet_gateway" "internet_gateway" {
  compartment_id = var.compartment_ocid
  vcn_id         = oci_core_vcn.gr_vcn.id
  display_name   = "${var.vcn_name}-internet-gateway"
}

resource "oci_core_route_table" "route_table" {
  compartment_id = var.compartment_ocid
  vcn_id         = oci_core_vcn.gr_vcn.id
  display_name   = "${var.vcn_name}-route-table"

  route_rules {
    destination       = "0.0.0.0/0"
    network_entity_id = oci_core_internet_gateway.internet_gateway.id
  }
}

resource "oci_core_security_list" "bastion_security_list" {
  compartment_id = var.compartment_ocid
  vcn_id         = oci_core_vcn.gr_vcn.id
  display_name   = "${var.vm_name}-security-list"

  # Allow inbound SSH traffic
  ingress_security_rules {
    protocol  = "6" # TCP
    source    = "0.0.0.0/0"
    stateless = false

    tcp_options {
      min = 22
      max = 22
    }
  }
  
  # Allow all outbound traffic
  egress_security_rules {
    protocol    = "all"
    destination = "0.0.0.0/0"
    stateless   = false
  }
}

resource "oci_core_security_list" "benchmark_security_list" {
  compartment_id = var.compartment_ocid
  vcn_id         = oci_core_vcn.gr_vcn.id
  display_name   = "${var.bm_name}-security-list"

  # Allow inbound SSH traffic
  ingress_security_rules {
    protocol  = "6" # TCP
    source    = "0.0.0.0/0"
    stateless = false

    tcp_options {
      min = 22
      max = 22
    }
  }
  
  # Allow inbound SSH traffic
  ingress_security_rules {
    protocol  = "6" # TCP
    source    = var.bastion_subnet_cidr
    stateless = false

    tcp_options {
      min = 6443
      max = 6443
    }
  }
  
  # Allow all outbound traffic
  egress_security_rules {
    protocol    = "all"
    destination = "0.0.0.0/0"
    stateless   = false
  }
}

resource "oci_core_subnet" "bastion_subnet" {
  availability_domain        = var.availability_domain
  cidr_block                 = var.bastion_subnet_cidr
  display_name               = "${var.vm_name}-subnet"
  dns_label                  = "vmsubnet"
  compartment_id             = var.compartment_ocid
  vcn_id                     = oci_core_vcn.gr_vcn.id
  route_table_id             = oci_core_route_table.route_table.id
  security_list_ids          = [oci_core_security_list.bastion_security_list.id]
  prohibit_public_ip_on_vnic = false
}

resource "oci_core_subnet" "benchmark_subnet" {
  availability_domain        = var.availability_domain
  cidr_block                 = var.benchmark_subnet_cidr
  display_name               = "${var.bm_name}-subnet"
  dns_label                  = "bmsubnet"
  compartment_id             = var.compartment_ocid
  vcn_id                     = oci_core_vcn.gr_vcn.id
  route_table_id             = oci_core_route_table.route_table.id
  security_list_ids          = [oci_core_security_list.benchmark_security_list.id]
  prohibit_public_ip_on_vnic = false
}

resource "oci_core_instance" "bastion_instance" {
  availability_domain = var.availability_domain
  compartment_id      = var.compartment_ocid
  display_name        = var.vm_name
  shape               = var.vm_shape

  shape_config {
    ocpus         = var.vm_ocpus
    memory_in_gbs = var.vm_memory_in_gbs
  }

  create_vnic_details {
    subnet_id        = oci_core_subnet.bastion_subnet.id
    hostname_label   = var.vm_name
    assign_public_ip = true
  }

  source_details {
    source_type = "image"
    source_id   = data.oci_core_images.ubuntu_images.images[0].id
  }

  metadata = {
    ssh_authorized_keys = var.ssh_public_key
  }

  # Prevent rapid recycling of the instance
  preserve_boot_volume = false

  connection {
    user        = var.ssh_user
    private_key = file(var.ssh_private_key_path)
    host        = self.public_ip
  }
}

resource "oci_core_instance" "benchmark_instance" {
  availability_domain = var.availability_domain
  compartment_id      = var.compartment_ocid
  display_name        = var.bm_name
  shape               = var.bm_shape

  depends_on = [oci_core_instance.bastion_instance]

  create_vnic_details {
    subnet_id        = oci_core_subnet.benchmark_subnet.id
    hostname_label   = var.bm_name
    assign_public_ip = true
  }

  source_details {
    source_type = "image"
    source_id   = data.oci_core_images.ubuntu_images.images[0].id
  }

  metadata = {
    ssh_authorized_keys = var.ssh_public_key
  }

  # Prevent rapid recycling of the instance
  preserve_boot_volume = false

  connection {
    user        = var.ssh_user
    private_key = file(var.ssh_private_key_path)
    host        = self.public_ip
  }

  # Ubuntu image from Oracle has default iptables rules that drop traffic which
  # breaks cluster networking and must be removed.
  provisioner "file" {
    source      = "${path.module}/scripts/clean_iptables.sh"
    destination = "/tmp/clean_iptables.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/clean_iptables.sh",
      "/tmp/clean_iptables.sh"
    ]
  }

  provisioner "remote-exec" {
    inline = [
      "curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION='${var.k3s_version}' sh -",
      "systemctl is-active --quiet k3s.service",
    ]
  }

  timeouts {
    create = "20m" # Bare metal instances can take longer to provision
  }
}

resource "null_resource" "fetch_kubeconfig" {
  depends_on = [oci_core_instance.benchmark_instance]
  triggers = {
    always_run = "${timestamp()}"
  }

  connection {
    user        = var.ssh_user
    private_key = file(var.ssh_private_key_path)
    host        = oci_core_instance.benchmark_instance.public_ip
  }

  provisioner "remote-exec" {
    inline = [
      "mkdir $HOME/.kube",
      "sudo cp /etc/rancher/k3s/k3s.yaml $HOME/.kube/config",
      "sudo chown ${var.ssh_user}:${var.ssh_user} $HOME/.kube/config",
      "export KUBECONFIG=$HOME/.kube/config",
    ]
  }

  provisioner "local-exec" {
    command = "scp -i ${var.ssh_private_key_path} -o StrictHostKeyChecking=no ${var.ssh_user}@${oci_core_instance.benchmark_instance.public_ip}:/home/${var.ssh_user}/.kube/config ./kube-config"
  }  
}

output "bastion_node_ip" {
  value = oci_core_instance.bastion_instance.public_ip
}

output "benchmark_node_ip" {
  value = oci_core_instance.benchmark_instance.public_ip
}
