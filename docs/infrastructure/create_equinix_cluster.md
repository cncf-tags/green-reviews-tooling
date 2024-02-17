# Create Equinix Metal cluster using CAPI / CAPEM

## Overview

- We are using machines contributed to the CNCF by Equinix Metal.
- These steps manually create a cluster using CAPI / CAPEM using a local
Kind cluster for the management cluster.
- These docs are based on the Equinix CAPI [guide](https://deploy.equinix.com/developers/guides/kubernetes-cluster-api/).
- Using [Podman](https://github.com/kubernetes-sigs/cluster-api-provider-packet/issues/496)
is recommended by Equinix for MacOS.
- We plan to automate these steps later using GitHub Actions and an IaC tool.

## Create management Cluster

- If you don't have a management cluster running create it with
```sh
kind create cluster
```

- Check that it's running locally

```sh
kubectl cluster-info --context kind-kind
```

## Create Cluster

- Log in to the Equinix Metal [console](https://console.equinix.com/) and create a new project. Get
the project API key from the project settings. If it doesn't exist create it.

```sh
export PACKET_API_KEY="<PROJECT_API_KEY>"
```

- Create a project ssh-key with `<YOUR_SSH_KEY>`

- Install the CAPI controllers using [clusterctl](https://cluster-api.sigs.k8s.io/user/quick-start.html#install-clusterctl).

```sh
clusterctl init --infrastructure packet
```

- Set env vars with cluster config.

```sh
# Get the project ID from the project settings in the console
export PROJECT_ID="<PROJECT_ID>"

# Use Paris metro (Equinix region)
export METRO="pa"

# Use Ubuntu 22.04 with cGroup v2 for Kepler
export NODE_OS="ubuntu_22_04"

# The pod and service CIDRs for the new cluster
export POD_CIDR="192.168.0.0/16"
export SERVICE_CIDR="172.26.0.0/16"

# Use node type with Intel CPUs for RAPL support
export CONTROLPLANE_NODE_TYPE="m3.small.x86"
export WORKER_NODE_TYPE="m3.small.x86"

# SSH key to use for access to nodes
export SSH_KEY="<YOUR_SSH_KEY>"

# Kubernetes version to install
export KUBERNETES_VERSION="v1.28.2"
```

- Generate cluster manifests.

```sh
clusterctl generate cluster wg-green-reviews \
  --kubernetes-version $KUBERNETES_VERSION \
  --control-plane-machine-count=1 \
  --worker-machine-count=1 \
  > wg-green-reviews.yaml
```

- Apply cluster manifests.

```sh
kubectl apply -f wg-green-reviews.yaml
```

- Get kubeconfig and store it securely.

```sh
clusterctl get kubeconfig wg-green-reviews > wg-green-reviews.kubeconfig
```

- Set `KUBECONFIG` env var so following commands are run on the cluster.

```sh
export KUBECONFIG=wg-green-reviews.kubeconfig
```

- Install Cilium as CNI.

```sh
helm repo add cilium https://helm.cilium.io/
helm install cilium cilium/cilium  --version 1.14.2 --namespace kube-system
```

- SSH to each cluster node and ensure Kepler dependencies are installed
(user is named `root`).

```sh
apt update
apt install -y linux-headers-$(uname -r)
apt install -y linux-modules-extra-$(uname -r) 
modprobe intel_rapl_common
```

- Install Kepler.

```sh
helm repo add kepler https://sustainable-computing-io.github.io/kepler-helm-chart
helm install kepler kepler/kepler --namespace kepler --create-namespace
```

- Check Kepler container metrics are non-zero. 

```sh
kubectl exec -ti -n kepler daemonset/kepler \
    -- bash -c "curl localhost:9102/metrics" | grep 'kepler_container_package_joules_total'
```

## Delete Cluster

- If the Kind cluster still exists it can be used to delete the Equinix cluster.

```sh
kubectl delete cluster wg-green-reviews
```

- Otherwise delete both servers and the elastic IP via the Equinix console.
