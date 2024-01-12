# Full Kepler install with dashboard


## Install `kube-prometheus-stack`
```bash
helm install prom prometheus-community/kube-prometheus-stack --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues
=false  --namespace monitoring --create-namespace
```

## Install `kepler`

> Make sure you override the specific fields
```diff
serviceMonitor:
-  enabled: false
+  enabled: true
-  namespace: ""
+  namespace: monitoring
```

```bash
helm install kepler kepler/kepler --values values.yaml -n monitoring
```

## Add the kepler grafana dashboard

```bash
cd clusters/
kubectl apply -f kepler-grafana.yaml
```
