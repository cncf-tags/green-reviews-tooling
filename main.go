// Dagger module for GreenReviewsTooling functions.

package main

import (
	"context"
	"os"

	"github.com/cncf-tags/green-reviews-tooling/internal/dagger"
	"github.com/cncf-tags/green-reviews-tooling/pkg/pipeline"
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
	p, err := newPipeline(source, kubeconfig)
	if err != nil {
		return nil, err
	}

	return p.Benchmark(ctx, cncfProject, config, version, benchmarkJobURL, benchmarkJobDurationMins)
}

func newPipeline(source *dagger.Directory, kubeconfig string) (*pipeline.Pipeline, error) {
	var configFile *dagger.File
	var err error

	container := build(source)
	configFile, err = getKubeconfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return pipeline.New(container, source, configFile)
}

func build(src *dagger.Directory) *dagger.Container {
	return dag.Container().
		WithDirectory("/src", src).
		Directory("/src").
		DockerBuild().
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")
}

func getKubeconfig(configFilePath string) (*dagger.File, error) {
	contents, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	filePath := "/.kube/config"
	dir := dag.Directory().WithNewFile(filePath, string(contents))
	return dir.File(filePath), nil
}
