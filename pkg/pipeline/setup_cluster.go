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
	_, err := p.exec(ctx, cmd.FluxInstall())
	if err != nil {
		return nil, err
	}

	// Add node labels.
	nodeName, err := p.getNodeName(ctx)
	if err != nil {
		return nil, err
	}
	_, err = p.exec(ctx, cmd.LabelNode(nodeName, nodeLabels()))
	if err != nil {
		return nil, err
	}

	// Apply cluster manifests.
	for _, manifest := range clusterManifests() {
		_, err = p.execWithDir(ctx, manifest, cmd.Apply(manifest))
		if err != nil {
			return nil, err
		}
	}

	// Patch helmrelease values to ensure all pods will start.
	for _, patch := range localSetupPatches() {
		_, err = p.exec(ctx, patch)
		if err != nil {
			return nil, err
		}
	}

	// Kepler depends on kube-prometheus-stack.
	_, err = p.exec(ctx, cmd.FluxReconcile("helmrelease", "kepler"))
	if err != nil {
		return nil, err
	}

	// Wait until all pods are ready.
	_, err = p.exec(ctx, cmd.WaitForNamespace(monitoringNamespace))
	if err != nil {
		return nil, err
	}

	// Enable terminal to debug.
	// return p.Terminal(ctx)

	return p.container, nil
}

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

func nodeLabels() map[string]string {
	return map[string]string{
		"cncf-project":                      "wg-green-reviews",
		"cncf-project-sub":                  "internal",
		"node-role.kubernetes.io/benchmark": "true",
	}
}
