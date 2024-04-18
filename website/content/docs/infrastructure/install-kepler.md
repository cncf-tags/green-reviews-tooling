---
title: "Install Kepler"
description: "This section provides details around installation of Kepler tool."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
slug: kepler-install
weight: 870
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

Full installation of Kepler with dashboard.

## Install kube-prometheus-stack

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false --namespace monitoring --create-namespace
```

## Install kepler

```bash
helm repo add kepler https://sustainable-computing-io.github.io/kepler-helm-chart
helm install kepler kepler/kepler --set serviceMonitor.enabled=true --set serviceMonitor.namespace=monitoring --namespace monitoring
```

- If Kepler pods won't start or metrics are 0 check [troubleshooting](https://sustainable-computing.io/usage/trouble_shooting) documentation.

## Add the kepler grafana dashboard

```bash
cd clusters
kubectl apply -f kepler-grafana.yaml
```
