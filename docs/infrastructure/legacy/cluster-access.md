---
title: "Cluster Access"
description: "This section provides details for how to access the infrastructure Kubernetes cluster."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
slug: cluster-access
weight: 890
toc: true
---

To access the cluster we generate multiple kubeconfigs with different permissions.
These are stored in the TAG ENV 1Password account.

Currently supported kubeconfigs are:

- `pipeline` for our benchmarking pipeline. Added as a GitHub secret to this repo
- `readonly` for contributors to the project

## `pipeline` kubeconfig

- View cluster role + full access to Flux custom resources

### Steps to get the `pipeline` kubeconfig

```bash
chmod u+x scripts/gen-kubeconfig.sh
./scripts/gen-kubeconfig.sh pipeline
```

## `readonly` kubeconfig

- View cluster role + port forwarding

> Refer: [Link](https://codeforphilly.github.io/chime/operations/limited-kubeconfigs/limited-kubeconfigs.html)
[create-token](https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/#create-token)

### Steps to get the `readonly` kubeconfig

```bash
chmod u+x scripts/gen-kubeconfig.sh
./scripts/gen-kubeconfig.sh readonly
```

## Test out the kubeconfig

```bash
export KUBECONFIG=${PWD}/green-reviews-cluster-readonly-config
```

```bash
kubectl get no # it will fail
```

## Access to resources

Each kubeconfig provides access to CRUD operations on different resources on the cluster.
Changes to the resources that can be accessed by each kubeconfig can be made by appending resources to the `ClusterRole` found in `./base/` such as, for example, `pipeline-kubeconfig-resources.yaml`.
Changes to a `ClusterRole` are reconciled on the cluster by Flux.
