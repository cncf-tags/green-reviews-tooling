# Local setup

1. Start Kubernetes
2. [Install and start Prometheus](https://sustainable-computing.io/installation/kepler/#deploy-the-prometheus-operator)
   1. `cd kube-prometheus`
   2. `kubectl apply --server-side -f manifests/setup`
   3. `kubectl apply -f manifests/`
   4. Wait…
   5. `kubectl -n monitoring port-forward svc/grafana 3000`
   6. Open dashboard _localhost:3000_
4. Install metrics server
   1. `kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml`
   2. Patch

```shell
kubectl patch -n kube-system deployment metrics-server --type=json \
-p '[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'
```

4. Install and start Kepler
   1. Open a new terminal
   2. `git clone --depth 1 git@github.com:sustainable-computing-io/kepler.git`
5. `cd kepler`
6. `make build-manifest OPTS="PROMETHEUS_DEPLOY"`
7. `kubectl apply -f _output/generated-manifest/deployment.yaml`
8. Add [dashboard](https://raw.githubusercontent.com/sustainable-computing-io/kepler/main/grafana-dashboards/Kepler-Exporter.json) to Grafana.
9. Install and start Falco
   1. Open a new terminal
   2. [Install Helm](https://helm.sh/docs/intro/install/)
   3. `helm repo add falcosecurity https://falcosecurity.github.io/charts`
   4. `helm repo update`
   5. `helm install falco falcosecurity/falco --namespace falco --create-namespace --set driver.kind=modern-bpf  --set falco.grpc.enabled=true --set falco.grpc_output.enabled=true`
   6. `helm install falco-exporter falcosecurity/falco-exporter`
10. Run Falco tests
    1. https://github.com/falcosecurity/cncf-green-review-testing/tree/main/benchmark-tests
    2. May need to remove `nodeSelector`
    3. Write out metrics to JSON
    4. Thinking about https://github.com/prometheus/prom2json
    5. These metrics:

      ```
      rate(container_cpu_usage_seconds_total[5m])
      container_memory_rss
      container_memory_working_set_bytes
      kepler_container_joules_total
      ```
