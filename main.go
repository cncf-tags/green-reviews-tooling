// Dagger module for GreenReviewsTooling functions.

package main

import (
	"context"
	"log"
	"os"

	"github.com/cncf-tags/green-reviews-tooling/internal/dagger"
	"github.com/cncf-tags/green-reviews-tooling/pkg/pipeline"
)

const (
	clusterName = "green-reviews-test"
	sourceDir   = "src"
)

type GreenReviewsTooling struct{}

// BenchmarkPipeline measures the sustainability footprint of CNCF projects.
func (m *GreenReviewsTooling) BenchmarkPipeline(ctx context.Context,
	source *dagger.Directory,
	cncfProject,
	// +optional
	config,
	version,
	benchmarkJobURL,
	kubeconfig string,
	prometheus_url string,
	benchmarkJobDurationMins int) (*dagger.Container, error) {
	p, err := newPipeline(ctx, source, kubeconfig)
	if err != nil {
		return nil, err
	}

	return p.Benchmark(ctx, cncfProject, config, version, benchmarkJobURL, benchmarkJobDurationMins, prometheus_url)
}

// BenchmarkPipelineTest tests the pipeline.
func (m *GreenReviewsTooling) BenchmarkPipelineTest(ctx context.Context,
	source *dagger.Directory,
	// +optional
	// +default="falco"
	cncfProject,
	// +optional
	// +default="modern-ebpf"
	config,
	// +optional
	// +default="0.40.0"
	version,
	// +optional
	// +default="https://raw.githubusercontent.com/falcosecurity/cncf-green-review-testing/2551137b1a09bd0594f76b09e82e08c98f95efd3/benchmark-tests/falco-benchmark-tests.yaml"
	benchmarkJobURL,
	// +optional
	kubeconfig string,
	// +optional
	prometheus_url string,
	// +optional
	// +default=2
	benchmarkJobDurationMins int) (*dagger.Container, error) {
	p, err := newPipeline(ctx, source, kubeconfig)
	if err != nil {
		return nil, err
	}

	if kubeconfig == "" {
		// This is a new k3s container so bootstrap flux and the monitoring stack.
		_, err = p.SetupCluster(ctx)
		if err != nil {
			return nil, err
		}
	}

	if prometheus_url == "" {
		log.Printf("Missing prometheus url from makefile", err)
		return nil, err
	}

	return p.Benchmark(ctx,
		cncfProject,
		config,
		version,
		benchmarkJobURL,
		benchmarkJobDurationMins,
		prometheus_url)
}

// SetupCluster installs cluster components in an empty cluster for CI/CD and
// local development.
func (m *GreenReviewsTooling) SetupCluster(ctx context.Context,
	source *dagger.Directory,
	// +optional
	kubeconfig string) (*dagger.Container, error) {
	p, err := newPipeline(ctx, source, kubeconfig)
	if err != nil {
		return nil, err
	}

	return p.SetupCluster(ctx)
}

// Terminal returns dagger interactive terminal configured with kubeconfig
// for trouble shooting.
func (m *GreenReviewsTooling) Terminal(ctx context.Context,
	source *dagger.Directory,
	// +optional
	kubeconfig string) (*dagger.Container, error) {
	p, err := newPipeline(ctx, source, kubeconfig)
	if err != nil {
		return nil, err
	}

	return p.Terminal(ctx)
}

// newPipeline creates a new benchmark pipeline.
func newPipeline(ctx context.Context, source *dagger.Directory, kubeconfig string) (*pipeline.Pipeline, error) {
	var configFile *dagger.File
	var err error

	container := build(source)

	if kubeconfig == "" {
		// No kubeconfig so start a new k3s container using the k3s dagger
		// module. See https://daggerverse.dev/mod/github.com/marcosnils/daggerverse/k3s
		configFile, err = startK3sCluster(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		// Connect to an existing cluster via a kubeconfig.
		configFile, err = getKubeconfig(kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	return pipeline.New(container, source, configFile)
}

// build builds a container image from the Dockerfile with any CLIs needed by
// the pipeline such as kubectl.
func build(src *dagger.Directory) *dagger.Container {
	return dag.Container().
		WithDirectory(sourceDir, src).
		Directory(sourceDir).
		DockerBuild().
		WithMountedDirectory(sourceDir, src).
		WithWorkdir(sourceDir)
}

// getKubeconfig returns a dagger file object pointing as the cluster kubeconfig.
func getKubeconfig(configFilePath string) (*dagger.File, error) {
	contents, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	filePath := pipeline.KubeconfigPath
	dir := dag.Directory().WithNewFile(filePath, string(contents))
	return dir.File(filePath), nil
}

// startK3sCluster starts a new k3s container using the k3s dagger module and
// returns its kubeconfig.
func startK3sCluster(ctx context.Context) (*dagger.File, error) {
	k3s := dag.K3S(clusterName)
	kServer := k3s.Server()
	if _, err := kServer.Start(ctx); err != nil {
		return nil, err
	}
	return k3s.Config(), nil
}
