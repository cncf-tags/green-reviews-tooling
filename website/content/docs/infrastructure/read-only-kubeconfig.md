---
title: "Cluster Access"
description: "This section provides details for how to access the infrastructure Kubernetes cluster."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
slug: cluster-access
weight: 860
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

## Read only kubeconfig

> Refer: [Link](https://codeforphilly.github.io/chime/operations/limited-kubeconfigs/limited-kubeconfigs.html)

[create-token](https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/#create-token)

## Steps to get the kubeconfig

```bash
chmod u+x scripts/gen-readonly-kubeconfig.sh
./scripts/gen-readonly-kubeconfig.sh
```

## Test out the kubeconfig

```bash
export KUBECONFIG=${PWD}/green-reviews-cluster-readonly-config
```

```bash
kubectl get no # it will fail
```

> Just get, watch, list for all pods
