// Dagger module for GreenReviewsTooling functions.

package main

import (
	"context"
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
	benchmarkJobDurationMins int) (*dagger.Container, error) {
	p, err := newPipeline(ctx, source, kubeconfig)
	if err != nil {
		return nil, err
	}

	return p.Benchmark(ctx, cncfProject, config, version, benchmarkJobURL, benchmarkJobDurationMins)
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
	// +default="0.39.2"
	version,
	// +optional
	// +default="https://raw.githubusercontent.com/falcosecurity/cncf-green-review-testing/e93136094735c1a52cbbef3d7e362839f26f4944/benchmark-tests/falco-benchmark-tests.yaml"
	benchmarkJobURL,
	// +optional
	kubeconfig string,
	// +optional
	// +default=2
	benchmarkJobDurationMins int) (*dagger.Container, error) {
	p, err := newPipeline(ctx, source, kubeconfig)
	if err != nil {
		return nil, err
	}

	if kubeconfig == "" {
		_, err = p.SetupCluster(ctx)
		if err != nil {
			return nil, err
		}
	}

	return p.Benchmark(ctx,
		cncfProject,
		config,
		version,
		benchmarkJobURL,
		benchmarkJobDurationMins)
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

func newPipeline(ctx context.Context, source *dagger.Directory, kubeconfig string) (*pipeline.Pipeline, error) {
	var configFile *dagger.File
	var err error

	container := build(source)

	if kubeconfig == "" {
		configFile, err = startK3sCluster(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		configFile, err = getKubeconfig(kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	return pipeline.New(container, source, configFile)
}

func build(src *dagger.Directory) *dagger.Container {
	return dag.Container().
		WithDirectory(sourceDir, src).
		Directory(sourceDir).
		DockerBuild().
		WithMountedDirectory(sourceDir, src).
		WithWorkdir(sourceDir)
}

func getKubeconfig(configFilePath string) (*dagger.File, error) {
	contents, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	filePath := pipeline.KubeconfigPath
	dir := dag.Directory().WithNewFile(filePath, string(contents))
	return dir.File(filePath), nil
}

func startK3sCluster(ctx context.Context) (*dagger.File, error) {
	k3s := dag.K3S(clusterName)
	kServer := k3s.Server()
	if _, err := kServer.Start(ctx); err != nil {
		return nil, err
	}
	return k3s.Config(), nil
}
