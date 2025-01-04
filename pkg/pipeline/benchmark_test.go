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
	benchmarkJobDurationMins int) (*dagger.Container, error) {
	_, err := p.Benchmark(ctx, cncfProject, config, version, benchmarkJobURL, benchmarkJobDurationMins)
	if err != nil {
		log.Printf("benchmark failed: %v", err)
	}

	// TODO Add tests.
	return p.container, nil
}
