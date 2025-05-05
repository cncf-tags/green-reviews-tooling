#!/bin/bash

# Get K3s version from environment variable
K3S_VERSION=${K3S_VERSION}

# Get the public IP using ifconfig.me
PUBLIC_IP=$(curl -s ifconfig.me)

curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION="${K3S_VERSION}" \
  K3S_KUBECONFIG_MODE="644" \
  INSTALL_K3S_EXEC="server \
    --tls-san ${PUBLIC_IP} \
    --bind-address 0.0.0.0 \
    --https-listen-port 6443 \
    --advertise-address ${PUBLIC_IP}" sh -
systemctl is-active --quiet k3s.service
