# Dagger

Run benchmark pipeline locally using [dagger](https://docs.dagger.io/).

Docs use `kind` other tooling like minikube or k3d should also work but is untested.

## Tools

These CLIs are needed

- `dagger` https://docs.dagger.io/install
- `helm` https://helm.sh/docs/helm/helm_install/
- `kind` https://kind.sigs.k8s.io/docs/user/quick-start/
- `yq` https://github.com/mikefarah/yq/#install

## Setup

- Create kind cluster

```sh
kind create cluster
kind get kubeconfig | yq e '.clusters[0].cluster.server = "https://kubernetes.default"' - > kind-in-cluster
```

- Install Dagger engine and configure CLI https://docs.dagger.io/integrations/kubernetes

```sh
helm upgrade --install --namespace=dagger --create-namespace \
    dagger oci://registry.dagger.io/dagger-helm

kubectl wait --for condition=Ready --timeout=60s pod \
    --selector=name=dagger-dagger-helm-engine --namespace=dagger

DAGGER_ENGINE_POD_NAME="$(kubectl get pod \
    --selector=name=dagger-dagger-helm-engine --namespace=dagger \
    --output=jsonpath='{.items[0].metadata.name}')"
export DAGGER_ENGINE_POD_NAME

_EXPERIMENTAL_DAGGER_RUNNER_HOST="kube-pod://$DAGGER_ENGINE_POD_NAME?namespace=dagger"
export _EXPERIMENTAL_DAGGER_RUNNER_HOST
```

## Run pipeline

- Bootstrap flux and install manifests from [/clusters/base/](/clusters/base/)

```sh
dagger call setup-cluster --source=. --kubeconfig=/src/kind-in-cluster
```

- Run the pipeline and execute tests on completion

```sh
dagger call benchmark-pipeline-test \
    --source=. --kubeconfig=/src/kind-in-cluster \
    --cncf-project='falco' \
    --config='modern-ebpf' \
    --version='0.39.2' \
    --benchmark-job-url='https://raw.githubusercontent.com/falcosecurity/cncf-green-review-testing/e93136094735c1a52cbbef3d7e362839f26f4944/benchmark-tests/falco-benchmark-tests.yaml' \
    --benchmark-job-duration-mins=2
```

## Debugging

- Get an interactive terminal for trouble shooting

```sh
dagger call terminal --source=. --kubeconfig=/src/kind-in-cluster
```
