# Default target
.PHONY: all
all: verify

KUBECONFIG = green-reviews-test-kubeconfig
BENCHMARK_OUTPUT = benchmark_results.json

# Get a terminal for debugging
.PHONY: debug
debug:
	dagger call terminal --source=. --kubeconfig=/src/$(KUBECONFIG)

# Regenerate client bindings for the Dagger API
.PHONY: develop
develop:
	dagger develop

# Install dagger
.PHONY: install
install:
	helm upgrade --install \
		--namespace=dagger \
		--create-namespace \
		--version 0.18.17 \
		dagger oci://registry.dagger.io/dagger-helm && \
	kubectl wait \
		--for condition=Ready \
		--timeout=120s pod \
		--selector=name=dagger-dagger-helm-engine \
		--namespace=dagger && \
	DAGGER_ENGINE_POD_NAME=$$(kubectl get pod \
			--selector=name=dagger-dagger-helm-engine \
			--namespace=dagger \
			--output=jsonpath='{.items[0].metadata.name}') && \
	_EXPERIMENTAL_DAGGER_RUNNER_HOST="kube-pod://$$DAGGER_ENGINE_POD_NAME?namespace=dagger" && \
	echo "-------- Dagger installed successfully --------" && \
	echo " Install complete - add env vars to your shell" && \
	echo "  export DAGGER_ENGINE_POD_NAME=\"$$DAGGER_ENGINE_POD_NAME\"" && \
	echo "  export _EXPERIMENTAL_DAGGER_RUNNER_HOST=\"$$_EXPERIMENTAL_DAGGER_RUNNER_HOST\""

# Bootstrap cluster with flux and monitoring stack
.PHONY: setup
setup:
	dagger call setup-cluster \
		--source=. --kubeconfig=/src/$(KUBECONFIG)

# Test pipeline with default values
.PHONY: test
test:
	dagger call benchmark-pipeline-test --source=. --kubeconfig=/src/$(KUBECONFIG) export --path=$(BENCHMARK_OUTPUT)

# Verify tools are installed
.PHONY: verify
verify:
	dagger version
	helm version
	kubectl version --client
	yq --version
