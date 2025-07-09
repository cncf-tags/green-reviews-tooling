package pipeline

import (
	"context"
	"log"

	"github.com/cncf-tags/green-reviews-tooling/internal/dagger"
)

// BenchmarkTest runs the pipeline and executes tests on completion.
func (p *Pipeline) BenchmarkTest(ctx context.Context,
	cncfProject,
	config,
	version,
	benchmarkJobURL string,
	benchmarkJobDurationMins int,
	prometheus_url string) (*dagger.Container, error) {
	if _, err := p.Benchmark(ctx, cncfProject, config, version, benchmarkJobURL, benchmarkJobDurationMins, prometheus_url); err != nil {
		log.Printf("benchmark failed: %v", err)
	}

	// TODO Add tests.
	return p.container, nil
}
