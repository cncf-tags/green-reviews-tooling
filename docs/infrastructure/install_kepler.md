# Install Kepler

Full installation of Kepler with dashboard.

## Install `kube-prometheus-stack`

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false --namespace monitoring --create-namespace
```

## Install `kepler`

```bash
helm repo add kepler https://sustainable-computing-io.github.io/kepler-helm-chart
helm install kepler kepler/kepler --set serviceMonitor.enabled=true --set serviceMonitor.namespace=monitoring --namespace monitoring
```

- If Kepler pods won't start or metrics are 0 check [trouble shooting](https://sustainable-computing.io/usage/trouble_shooting/) docs.

## Add the kepler grafana dashboard

```bash
cd clusters/
kubectl apply -f kepler-grafana.yaml
```
