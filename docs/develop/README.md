# Development

We use [dagger](https://docs.dagger.io/) and its Go SDK to run the pipeline
locally during development and in automated tests using GitHub Actions.

## Tools

Docs use [kind](https://kind.sigs.k8s.io/docs/user/quick-start/) to create a local cluster. 
Other tooling like `minikube` or `k3d` should also work but is untested.

These additional tools need to be installed.

- `dagger` https://docs.dagger.io/install
- `kubectl` https://kubernetes.io/docs/tasks/tools/
- `helm` https://helm.sh/docs/helm/helm_install/
- `yq` https://github.com/mikefarah/yq/#install

## Setup

- Verify CLIs are installed

```sh
make
```

- Create kind cluster and add kubeconfig to source dir so dagger can access it.

```sh
kind create cluster
kind get kubeconfig | yq e '.clusters[0].cluster.server = "https://kubernetes.default"' - > green-reviews-test-kubeconfig
```

- Install dagger engine. 

```sh
make install
```

- Add `DAGGER_ENGINE_POD_NAME` and `_EXPERIMENTAL_DAGGER_RUNNER_HOST` env vars to your shell. See https://docs.dagger.io/integrations/kubernetes/#example

- Bootstrap cluster with flux and monitoring stack.

```sh
make setup
```

## Test

- Run integration test.

```sh
make test
```

## Debugging

- Get an [interactive terminal](https://docs.dagger.io/api/terminal/) for trouble shooting.

```sh
make debug
```

## Development

- Run [dagger develop](https://docs.dagger.io/reference/cli/#dagger-develop) to
regenerate client bindings for the dagger API.

```sh
make develop
```
