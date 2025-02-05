package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/cncf-tags/green-reviews-tooling/internal/dagger"
	"github.com/cncf-tags/green-reviews-tooling/pkg/cmd"
)

const (
	monitoringNamespace = "monitoring"
)

// SetupCluster bootstraps flux and installs manifests from /clusters/base/ for
// CI/CD and local development.
func (p *Pipeline) SetupCluster(ctx context.Context) (*dagger.Container, error) {
	// Install flux.
	if _, err := p.exec(ctx, cmd.FluxInstall()); err != nil {
		return nil, err
	}

	// Add node labels.
	nodeName, err := p.getNodeName(ctx)
	if err != nil {
		return nil, err
	}
	if _, err = p.exec(ctx, cmd.LabelNode(nodeName, nodeLabels())); err != nil {
		return nil, err
	}

	// Apply cluster manifests.
	for _, manifest := range clusterManifests() {
		if _, err = p.execWithDir(ctx, manifest, cmd.Apply(manifest)); err != nil {
			return nil, err
		}
	}

	// Patch helmrelease values to ensure all pods will start.
	for _, patch := range localSetupPatches() {
		if _, err = p.exec(ctx, patch); err != nil {
			return nil, err
		}
	}

	// Kepler depends on kube-prometheus-stack so we wait till its reconciled.
	if _, err = p.exec(ctx, cmd.FluxReconcile("helmrelease", "kepler")); err != nil {
		return nil, err
	}

	// Wait until all pods are ready.
	if _, err = p.exec(ctx, cmd.WaitForReadyPods(monitoringNamespace)); err != nil {
		return nil, err
	}

	// Enable terminal to debug.
	// return p.Terminal(ctx)

	return p.container, nil
}

// getNodeName gets the name of the cluster node. In CI all clusters are single
// node.
func (p *Pipeline) getNodeName(ctx context.Context) (string, error) {
	stdout, err := p.exec(ctx, cmd.GetNodeNames())
	if err != nil {
		return "", err
	}

	parts := strings.Split(stdout, "\n")
	if len(parts) == 0 {
		return "", fmt.Errorf("failed to get node name from %s", stdout)
	}

	return parts[0], nil
}

// clusterManifests are applied to bootstrap the monitoring stack.
func clusterManifests() []string {
	return []string{
		// Namespace must be created first for dependencies.
		"/clusters/base/monitoring-namespace.yaml",
		"/clusters/base",
	}
}

// localSetupPatches patch flux manifests to disable resources not available
// in the k3s container like host mounts.
func localSetupPatches() [][]string {
	return [][]string{
		cmd.Patch("helmrelease",
			"kube-prometheus-stack",
			"flux-system",
			"/spec/values/prometheus-node-exporter",
			`{"hostRootFsMount": {"enabled": false}}`),
		cmd.Patch("helmrelease",
			"kepler",
			"flux-system",
			"/spec/values/canMount",
			`{"usrSrc": false}`),
	}
}

// nodeLabels are added to ensure they match label selectors in k8s manifests.
func nodeLabels() map[string]string {
	return map[string]string{
		"cncf-project":                      "wg-green-reviews",
		"cncf-project-sub":                  "internal",
		"node-role.kubernetes.io/benchmark": "true",
	}
}
